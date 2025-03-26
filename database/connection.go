// Package database handles database migrations and connections.
//
// This package includes functions for connecting to MySQL, running migrations,
// and handling database errors.
package database

import (
	"database/sql"
	"discord-bot-tickets/config"
	"errors"
	"fmt"
	"github.com/pterm/pterm"
	"log"
	"os"
	"os/exec"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// getDSN returns the Data Source Name (DSN) for the MySQL connection.
//
// It takes 2 arguments:
// - cfg: a pointer to the config.Config struct
// - withDB: a boolean value to determine if the database name should be included
//
// returns: a string value representing the DSN
func getDSN(cfg *config.Config, withDB bool) string {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/",
		cfg.DB.Username,
		cfg.DB.Password,
		cfg.DB.Host,
		cfg.DB.Port,
	)

	if withDB {
		dsn += fmt.Sprintf("%s?charset=utf8mb4&parseTime=True&loc=Local", cfg.DB.Table)
	} else {
		dsn += "?charset=utf8mb4&parseTime=True&loc=Local"
	}

	return dsn
}

// dropAllTables Drops all tables in the database.
//
// It takes 2 arguments:
// - db: a pointer to the sql.DB struct
// - dbName: a string value representing the database name
//
// returns: an error if any
func dropAllTables(db *sql.DB, dbName string) error {
	rows, err := db.Query(fmt.Sprintf("SELECT table_name FROM information_schema.tables WHERE table_schema = '%s'", dbName))
	if err != nil {
		return fmt.Errorf("failed to query tables: %v", err)
	}
	defer rows.Close()

	var tableName string
	for rows.Next() {
		if err := rows.Scan(&tableName); err != nil {
			return fmt.Errorf("failed to scan table name: %v", err)
		}
		_, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", tableName))
		if err != nil {
			return fmt.Errorf("failed to drop table %s: %v", tableName, err)
		}
		log.Printf("Table %s dropped successfully.", tableName)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating over tables: %v", err)
	}

	return nil
}

// MigrateDatabase Runs migrations on the MySQL database.
// If the --fresh flag is passed, all tables are dropped before running migrations.
//
// It takes 2 arguments:
// - cfg: a pointer to the config.Config struct
// - fresh: a boolean value to determine if all tables should be dropped
//
// returns: none
func MigrateDatabase(cfg *config.Config, fresh bool) {
	dsn := getDSN(cfg, true)

	// Open MySQL connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close MySQL connection: %v", err)
		}
	}(db)

	// if the --fresh flag is passed, drop all tables
	if fresh {
		err := dropAllTables(db, cfg.DB.Table)
		if err != nil {
			log.Fatalf("Failed to drop all tables: %v", err)
		}
	}

	// Run migrations
	m, err := migrate.New(
		"file://./database/migrations", // Path to your migration files with scheme
		fmt.Sprintf("mysql://%s", dsn), // MySQL connection string
	)

	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Apply migrations
	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("No changes in migration.")
		} else if migrateErr, ok := err.(migrate.ErrDirty); ok {
			result, _ := pterm.DefaultInteractiveConfirm.Show("Database is dirty. Would you like to force the migration?")

			pterm.Println()

			if !result {
				log.Fatalf("Migration failed: %v", migrateErr)
			}

			// get the current version
			version, _, _ := m.Version()

			// force the migration
			if err := m.Force(int(version)); err != nil {
				log.Fatalf("Migration failed: %v", err)
			}

			log.Println("Migration forced successfully!")
		} else {
			log.Fatalf("Migration failed: %v", err)
		}
	}
}

// Connect connects to the MySQL database using the provided configuration.
//
// It takes 1 argument:
// - cfg: a pointer to the config.Config struct
//
// returns: a pointer to the sql.DB struct and an error if any
func Connect(cfg *config.Config) (*sql.DB, error) {
	dsn := getDSN(cfg, true)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
		return nil, err
	}

	// Check if the connection is alive
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping MySQL: %v", err)
		return nil, err
	}

	return db, nil
}

// CreateMigration generates a new migration file with the given name.
//
// Parameters:
// - name: the name of the migration (e.g., "create_users_table").
// - dir: the path to the migrations folder (e.g., "db/migrations").
// - ext: the file extension (default is "sql").
//
// It automatically generates a sequential filename and creates the migration.
func CreateMigration(name string) {
	if name == "" {
		log.Fatalf("Migration name cannot be empty.")
	}

	// Sanitize the name (replace spaces with underscores)
	migrationName := strings.ReplaceAll(name, " ", "_")

	// Build the command
	cmd := exec.Command("migrate create -ext sql -dir database/migrations -seq ", migrationName)

	// Set output to terminal
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Execute the command
	log.Printf("Running: %s", strings.Join(cmd.Args, " "))
	if err := cmd.Run(); err != nil {
		log.Fatalf("Failed to create migration: %v", err)
	}

	log.Printf("Migration %s created successfully!", migrationName)
}

// RollbackMigration rolls back the last migration applied to the database.
//
// It takes 1 argument:
// - cfg: a pointer to the config.Config struct
//
// returns: none
func RollbackMigration(cfg *config.Config) {
	dsn := getDSN(cfg, true)

	// Open MySQL connection
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to MySQL: %v", err)
	}
	defer func(db *sql.DB) {
		err := db.Close()
		if err != nil {
			log.Fatalf("Failed to close MySQL connection: %v", err)
		}
	}(db)

	// Run migrations
	m, err := migrate.New(
		"file://./database/migrations",
		fmt.Sprintf("mysql://%s", dsn),
	)

	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	// Rollback the last migration
	if err := m.Steps(-1); err != nil {
		log.Fatalf("Rollback failed: %v", err)
	}

	log.Println("Migration rolled back successfully!")
}

package logger

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/natefinch/lumberjack.v2"
)

var Logger *log.Logger

func Init() error {
	// Create logs directory if it doesn't exist
	logsDir := "logs"
	if err := os.MkdirAll(logsDir, 0755); err != nil {
		return fmt.Errorf("failed to create logs directory: %v", err)
	}

	// Create log file with today's date
	currentDate := time.Now().Format("2006-01-02")
	logPath := filepath.Join(logsDir, fmt.Sprintf("%s.log", currentDate))

	// Setup lumberjack for log rotation
	lumberjackLogger := &lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    10,   // megabytes
		MaxBackups: 30,   // number of backups
		MaxAge:     30,   // days
		Compress:   true, // compress old files
	}

	// Create new logger
	Logger = log.New(lumberjackLogger, "", log.Ldate|log.Ltime|log.Lshortfile)

	return nil
}

// Info logs a message with INFO severity
func Info(format string, v ...interface{}) {
	Logger.Printf("[INFO] "+format, v...)
}

// Error logs a message with ERROR severity
func Error(format string, v ...interface{}) {
	Logger.Printf("[ERROR] "+format, v...)
}

// Warn logs a message with WARN severity
func Warn(format string, v ...interface{}) {
	Logger.Printf("[WARN] "+format, v...)
}

// Debug logs a message with DEBUG severity
func Debug(format string, v ...interface{}) {
	Logger.Printf("[DEBUG] "+format, v...)
}

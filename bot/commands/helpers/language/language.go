package language

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"sync"

	logger "discord-bot-tickets/logging"

	"github.com/diamondburned/arikawa/v3/discord"
)

// Translation represents a single translation entry
type Translation struct {
	Message string `json:"message"`
}

// LanguageFile represents the structure of our language JSON files
type LanguageFile struct {
	General struct {
		Errors struct {
			Channel    Translation `json:"channel"`
			Owner      Translation `json:"owner"`
			Generic    Translation `json:"generic"`
			NotATicket Translation `json:"not_a_ticket"`
			NoMessage  Translation `json:"no_message"`
		} `json:"errors"`
		Success struct {
			Generic Translation `json:"generic"`
		} `json:"success"`
	} `json:"general"`
	Commands struct {
		Close struct {
			Success Translation `json:"success"`
			Error   Translation `json:"error"`
		} `json:"close"`
		Reply struct {
			Success Translation `json:"success"`
			Error   Translation `json:"error"`
		} `json:"reply"`
	} `json:"commands"`
	Embeds struct {
		TicketClosed struct {
			Title       Translation `json:"title"`
			Description Translation `json:"description"`
			Footer      Translation `json:"footer"`
		} `json:"ticket_closed"`
	} `json:"embeds"`
}

var (
	translations = make(map[discord.Language]LanguageFile)
	defaultLang  = discord.EnglishUK
	selectedLang = defaultLang
	mu           sync.RWMutex
)

// LoadLanguage loads a language file from the specified path
func LoadLanguage(lang discord.Language, path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var langFile LanguageFile
	if err := json.Unmarshal(data, &langFile); err != nil {
		return err
	}

	mu.Lock()
	translations[lang] = langFile
	mu.Unlock()

	return nil
}

// LoadLanguages loads all language files from the specified directory
func LoadLanguages(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		lang := discord.Language(file.Name()[:len(file.Name())-5]) // Remove .json extension
		if err := LoadLanguage(lang, filepath.Join(dir, file.Name())); err != nil {
			return err
		}
	}

	return nil
}

// InitializeLanguage initializes the language system with the specified language
func InitializeLanguage(dir string) error {
	// Load all language files
	if err := LoadLanguages(dir); err != nil {
		return err
	}

	// Get the language from environment variable
	envLang := os.Getenv("BOT_LANGUAGE")
	if envLang == "" {
		envLang = "en-GB" // Default to English if not specified
	}

	// Check if the specified language exists
	lang := discord.Language(envLang)
	if _, exists := translations[lang]; !exists {
		// If the specified language doesn't exist, try to fall back to English
		if _, exists := translations[defaultLang]; !exists {
			logger.Error("Neither the specified language nor the default language (en-GB) exists!")
			panic("Language initialization failed")
		}
		logger.Error("Specified language " + envLang + " not found, falling back to en-GB")
		lang = defaultLang
	}

	// Set the selected language
	mu.Lock()
	selectedLang = lang
	mu.Unlock()

	logger.Info("Language system initialized with: " + string(selectedLang))
	return nil
}

// GetTranslation gets a translation for a specific key
func GetTranslation(key string) string {
	mu.RLock()
	defer mu.RUnlock()

	// Split the key by dots to navigate the structure
	// Example: "general.errors.channel" or "commands.close.success"
	parts := strings.Split(key, ".")

	// Get the translation based on the key
	var translation Translation
	switch parts[0] {
	case "general":
		switch parts[1] {
		case "errors":
			switch parts[2] {
			case "channel":
				translation = translations[selectedLang].General.Errors.Channel
			case "owner":
				translation = translations[selectedLang].General.Errors.Owner
			case "generic":
				translation = translations[selectedLang].General.Errors.Generic
			case "not_a_ticket":
				translation = translations[selectedLang].General.Errors.NotATicket
			case "no_message":
				translation = translations[selectedLang].General.Errors.NoMessage
			}
		case "success":
			switch parts[2] {
			case "generic":
				translation = translations[selectedLang].General.Success.Generic
			}
		}
	case "commands":
		switch parts[1] {
		case "close":
			switch parts[2] {
			case "success":
				translation = translations[selectedLang].Commands.Close.Success
			case "error":
				translation = translations[selectedLang].Commands.Close.Error
			}
		case "reply":
			switch parts[2] {
			case "success":
				translation = translations[selectedLang].Commands.Reply.Success
			case "error":
				translation = translations[selectedLang].Commands.Reply.Error
			}
		}
	case "embeds":
		switch parts[1] {
		case "ticket_closed":
			switch parts[2] {
			case "title":
				translation = translations[selectedLang].Embeds.TicketClosed.Title
			case "description":
				translation = translations[selectedLang].Embeds.TicketClosed.Description
			case "footer":
				translation = translations[selectedLang].Embeds.TicketClosed.Footer
			}
		}
	}

	return translation.Message
}

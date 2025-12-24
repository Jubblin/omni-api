package i18n

import (
	"embed"
	"encoding/json"
	"fmt"
	"sync"
)

//go:embed locales/*.json
var localesFS embed.FS

var (
	translations map[string]map[string]string
	currentLang  string
	mu           sync.RWMutex
)

// Init initializes the i18n system with the given language code
func Init(langCode string) error {
	mu.Lock()
	defer mu.Unlock()

	translations = make(map[string]map[string]string)
	return SetLanguage(langCode)
}

// SetLanguage sets the current language
func SetLanguage(langCode string) error {
	mu.Lock()
	defer mu.Unlock()

	// Load translation file
	data, err := localesFS.ReadFile(fmt.Sprintf("locales/%s.json", langCode))
	if err != nil {
		return fmt.Errorf("failed to load locale file for %s: %w", langCode, err)
	}

	var langTranslations map[string]string
	if err := json.Unmarshal(data, &langTranslations); err != nil {
		return fmt.Errorf("failed to parse locale file for %s: %w", langCode, err)
	}

	translations[langCode] = langTranslations
	currentLang = langCode
	return nil
}

// T translates a key with optional arguments
func T(key string, args ...interface{}) string {
	mu.RLock()
	defer mu.RUnlock()

	langTranslations, ok := translations[currentLang]
	if !ok {
		return key
	}

	translation, ok := langTranslations[key]
	if !ok {
		return key
	}

	if len(args) > 0 {
		return fmt.Sprintf(translation, args...)
	}

	return translation
}

// GetCurrentLanguage returns the current language code
func GetCurrentLanguage() string {
	mu.RLock()
	defer mu.RUnlock()
	return currentLang
}

// GetLanguageName returns the human-readable name for a language code
func GetLanguageName(langCode string) string {
	names := map[string]string{
		"en": "English",
		"es": "Spanish",
		"fr": "French",
		"de": "German",
	}
	if name, ok := names[langCode]; ok {
		return name
	}
	return langCode
}

// GetSupportedLanguages returns a list of supported language codes
func GetSupportedLanguages() []string {
	return []string{"en", "es", "fr", "de"}
}

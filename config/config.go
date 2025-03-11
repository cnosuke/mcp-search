package config

import (
	"github.com/cockroachdb/errors"
	"github.com/jinzhu/configor"
)

// Config represents the server configuration
type Config struct {
	Search SearchConfig `yaml:"search"`
}

// SearchConfig represents the search configuration
type SearchConfig struct {
	APIKey            string `yaml:"api_key" env:"BRAVE_SEARCH_API_KEY"`
	Timeout           int    `yaml:"timeout" env:"BRAVE_SEARCH_TIMEOUT" default:"30"`
	MaxRetries        int    `yaml:"max_retries" env:"BRAVE_SEARCH_MAX_RETRIES" default:"2"`
	DefaultCountry    string `yaml:"default_country" env:"BRAVE_SEARCH_COUNTRY" default:"US"`
	DefaultSearchLang string `yaml:"default_search_lang" env:"BRAVE_SEARCH_LANGUAGE" default:"en"`
	DefaultUILang     string `yaml:"default_ui_lang" env:"BRAVE_SEARCH_UI_LANGUAGE" default:"en-US"`
}

// LoadConfig loads the configuration from a file
func LoadConfig(path string) (*Config, error) {
	cfg := &Config{}
	if err := configor.New(&configor.Config{
		Debug:      false,
		Verbose:    false,
		Silent:     true,
		AutoReload: false,
	}).Load(cfg, path); err != nil {
		return nil, errors.Wrap(err, "failed to load configuration")
	}

	// Validate required fields
	if cfg.Search.APIKey == "" {
		return nil, errors.New("Brave Search API key is required (set in config.yml or BRAVE_SEARCH_API_KEY env var)")
	}

	return cfg, nil
}

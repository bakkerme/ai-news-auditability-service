package internal

import (
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

type Specification struct {
	RunDataTTLHours    int      `mapstructure:"RUN_DATA_TTL_HOURS"`
	CORSAllowedOrigins []string `mapstructure:"CORS_ALLOWED_ORIGINS"`
	LlmURL             string   `mapstructure:"LLM_URL"`
	LlmAPIKey          string   `mapstructure:"LLM_API_KEY"`
	LlmModel           string   `mapstructure:"LLM_MODEL"`
}

// Validate checks if the specification is valid
func (s *Specification) Validate() error {
	if s.RunDataTTLHours <= 0 {
		return fmt.Errorf("RunDataTTLHours must be positive")
	}
	return nil
}

// GetConfig loads the configuration from environment variables and .env file
func GetConfig() (*Specification, error) {
	v := viper.New()

	// Set default values
	v.SetDefault("RUN_DATA_TTL_HOURS", 168) // 7 days
	v.SetDefault("CORS_ALLOWED_ORIGINS", []string{"http://localhost:3000"})
	v.SetDefault("LLM_URL", "")
	v.SetDefault("LLM_API_KEY", "")
	v.SetDefault("LLM_MODEL", "gpt-4")

	// Configure Viper to read from .env file
	v.SetConfigName(".env") // Name of config file (without extension)
	v.SetConfigType("env")  // File type
	v.AddConfigPath(".")    // Path to look for the config file in

	// Attempt to read the config file, but don't fail if it's not found
	// This allows environment variables to take precedence or be the sole source
	_ = v.ReadInConfig()

	// Configure Viper to automatically read environment variables
	v.AutomaticEnv()
	// ANAS_RUN_DATA_TTL_HOURS will be read as RUN_DATA_TTL_HOURS
	v.SetEnvPrefix("anas") // Will look for env variables with "ANAS_" prefix
	// Replace dots and other characters with underscores for env var compatibility
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_", "-", "_"))

	var s Specification
	err := v.Unmarshal(&s)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal specification: %w", err)
	}

	if err := s.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}
	return &s, nil
}

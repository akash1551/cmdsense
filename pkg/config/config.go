package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
)

type Config struct {
	Provider      string `mapstructure:"provider"`   // "openai" or "ollama"
	OpenAIKey     string `mapstructure:"openai_key"` // OpenAI API Key
	OllamaBaseURL string `mapstructure:"ollama_base_url" default:"http://localhost:11434"`
	Model         string `mapstructure:"model"` // "gpt-4o", "llama3", etc.
}

var AppConfig Config

func LoadConfig() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(home, ".config", "cmdsense")
	err = os.MkdirAll(configPath, 0755)
	if err != nil {
		return err
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// Defaults
	viper.SetDefault("provider", "openai")
	viper.SetDefault("ollama_base_url", "http://localhost:11434")
	viper.SetDefault("model", "gpt-4o")

	// Env vars
	viper.SetEnvPrefix("CMDSENSE")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			// or create a default one
		} else {
			return fmt.Errorf("error reading config file: %w", err)
		}
	}

	err = viper.Unmarshal(&AppConfig)
	return err
}

func SaveConfig() error {
	if err := viper.SafeWriteConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileAlreadyExistsError); ok {
			return viper.WriteConfig()
		}
		return err
	}
	return nil
}

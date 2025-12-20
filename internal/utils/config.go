package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/austiecodes/goa/internal/consts"
	"github.com/austiecodes/goa/internal/types"
	"github.com/openai/openai-go/v3"
)

// OpenAIProviderConfig represents the OpenAI provider configuration
type OpenAIProviderConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url,omitempty"`
}

// ProviderConfigs holds all provider configurations
type ProviderConfigs struct {
	OpenAI OpenAIProviderConfig `json:"openai"`
}

// ModelConfig represents the model section in config
type ModelConfig struct {
	ChatModel  *types.Model `json:"chat_model,omitempty"`
	TitleModel *types.Model `json:"title_model,omitempty"`
	ThinkModel *types.Model `json:"think_model,omitempty"`
}

// Config represents the application configuration
type Config struct {
	Providers ProviderConfigs `json:"providers"`
	Model     ModelConfig     `json:"model"`
	Debug     bool            `json:"debug,omitempty"`
}

// DefaultConfig returns the default configuration
func DefaultConfig() *Config {
	return &Config{
		Providers: ProviderConfigs{
			OpenAI: OpenAIProviderConfig{},
		},
		Model: ModelConfig{
			ChatModel: &types.Model{
				Provider: consts.ProviderOpenAI,
				ModelID:  string(openai.ChatModelGPT5Nano),
			},
			TitleModel: &types.Model{
				Provider: consts.ProviderOpenAI,
				ModelID:  string(openai.ChatModelGPT5Nano),
			},
			ThinkModel: &types.Model{
				Provider: consts.ProviderOpenAI,
				ModelID:  string(openai.ChatModelGPT5Nano),
			},
		},
		Debug: false,
	}
}

// GetConfigPath returns the path to the configuration file
func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %v", err)
	}
	goaDir := filepath.Join(homeDir, consts.GoaDir)
	if err := os.MkdirAll(goaDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create goa directory: %v", err)
	}
	return filepath.Join(goaDir, ".goa"), nil
}

// LoadConfig loads the configuration from file
func LoadConfig() (*Config, error) {
	configPath, err := GetConfigPath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return default config
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	// Apply defaults for missing fields
	applyDefaults(&config)

	return &config, nil
}

// applyDefaults fills in default values for missing config fields
func applyDefaults(config *Config) {
	defaultConfig := DefaultConfig()

	// Apply default models if not set
	if config.Model.ChatModel == nil {
		config.Model.ChatModel = defaultConfig.Model.ChatModel
	}
	if config.Model.TitleModel == nil {
		config.Model.TitleModel = defaultConfig.Model.TitleModel
	}
	if config.Model.ThinkModel == nil {
		config.Model.ThinkModel = defaultConfig.Model.ThinkModel
	}
}

// SaveConfig saves the configuration to file
func SaveConfig(config *Config) error {
	configPath, err := GetConfigPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// GetOpenAIConfig returns the OpenAI provider configuration with defaults applied
func GetOpenAIConfig() (*OpenAIProviderConfig, error) {
	config, err := LoadConfig()
	if err != nil {
		return nil, err
	}

	openaiConfig := config.Providers.OpenAI

	// Apply default base URL if not set
	if openaiConfig.BaseURL == "" {
		openaiConfig.BaseURL = consts.DefaultBaseURL
	}

	return &openaiConfig, nil
}

// GetDebugMode returns whether debug mode is enabled
func GetDebugMode() bool {
	config, err := LoadConfig()
	if err != nil {
		return false
	}
	return config.Debug
}

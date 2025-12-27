package provider

import (
	"fmt"

	"github.com/austiecodes/goa/internal/client"
	"github.com/austiecodes/goa/internal/consts"
	anthropicprov "github.com/austiecodes/goa/internal/provider/anthropic"
	googleprov "github.com/austiecodes/goa/internal/provider/google"
	openaiprov "github.com/austiecodes/goa/internal/provider/openai"
	"github.com/austiecodes/goa/internal/utils"
)

func NewQueryClient(cfg *utils.Config, providerName string) (client.QueryClient, error) {
	switch providerName {
	case consts.ProviderOpenAI:
		openaiCfg := cfg.Providers.OpenAI
		if openaiCfg.APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key not configured. Please configure provider first")
		}
		baseURL := openaiCfg.BaseURL
		if baseURL == "" {
			baseURL = consts.DefaultBaseURL
		}
		return openaiprov.NewQueryClient(openaiCfg.APIKey, baseURL), nil
	case consts.ProviderGoogle:
		googleCfg := cfg.Providers.Google
		if googleCfg.APIKey == "" {
			return nil, fmt.Errorf("Google API key not configured. Please configure provider first")
		}
		return googleprov.NewQueryClient(googleCfg.APIKey, googleCfg.BaseURL), nil
	case consts.ProviderAnthropic:
		anthropicCfg := cfg.Providers.Anthropic
		if anthropicCfg.APIKey == "" {
			return nil, fmt.Errorf("Anthropic API key not configured. Please configure provider first")
		}
		// Anthropic SDK handles base URL internally via options if provided.
		return anthropicprov.NewQueryClient(anthropicCfg.APIKey, anthropicCfg.BaseURL), nil

	default:
		return nil, fmt.Errorf("unsupported provider: %s", providerName)
	}
}

// NewEmbeddingClient creates an embedding client for the specified provider.
func NewEmbeddingClient(cfg *utils.Config, providerName string) (client.EmbeddingClient, error) {
	switch providerName {
	case consts.ProviderOpenAI:
		openaiCfg := cfg.Providers.OpenAI
		if openaiCfg.APIKey == "" {
			return nil, fmt.Errorf("OpenAI API key not configured. Please configure provider first")
		}
		baseURL := openaiCfg.BaseURL
		if baseURL == "" {
			baseURL = consts.DefaultBaseURL
		}
		return openaiprov.NewEmbeddingClient(openaiCfg.APIKey, baseURL), nil
	case consts.ProviderGoogle:
		googleCfg := cfg.Providers.Google
		if googleCfg.APIKey == "" {
			return nil, fmt.Errorf("Google API key not configured. Please configure provider first")
		}
		return googleprov.NewEmbeddingClient(googleCfg.APIKey, googleCfg.BaseURL), nil
	// Anthropic doesn't support embeddings officially in the same way or requested yet.

	default:
		return nil, fmt.Errorf("unsupported embedding provider: %s", providerName)
	}
}

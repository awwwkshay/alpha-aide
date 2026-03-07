package agent

import (
	"fmt"

	"github.com/awwwkshay/alpha-aide/agent/config"
	anthropicprovider "github.com/awwwkshay/alpha-aide/llm/anthropic"
	geminiprovider "github.com/awwwkshay/alpha-aide/llm/gemini"
	openaiprovider "github.com/awwwkshay/alpha-aide/llm/openai"
)

// NewProvider constructs a Provider for the given provider name and model.
func NewProvider(cfg *config.Config, providerName, modelID string) (Provider, error) {
	switch providerName {
	case "anthropic":
		if cfg.AnthropicAPIKey == "" {
			return nil, fmt.Errorf("ANTHROPIC_API_KEY not set")
		}
		return anthropicprovider.New(cfg.AnthropicAPIKey, modelID, cfg.MaxTokens), nil
	case "openai":
		if cfg.OpenAIAPIKey == "" {
			return nil, fmt.Errorf("OPENAI_API_KEY not set")
		}
		return openaiprovider.New(cfg.OpenAIAPIKey, cfg.BaseURL, modelID, cfg.MaxTokens), nil
	case "gemini":
		if cfg.GeminiAPIKey == "" {
			return nil, fmt.Errorf("GEMINI_API_KEY not set")
		}
		return geminiprovider.New(cfg.GeminiAPIKey, modelID, cfg.MaxTokens), nil
	default:
		return nil, fmt.Errorf("unknown provider %q", providerName)
	}
}

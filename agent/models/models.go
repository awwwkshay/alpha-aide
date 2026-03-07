// Package models provides a catalog of known LLM models and providers.
package models

// Model describes a known LLM model.
type Model struct {
	Provider string
	ID       string
	Display  string
	Short    string
}

// All is the catalog of known models, grouped by provider.
var All = []Model{
	// Anthropic
	{"anthropic", "claude-opus-4-6", "Claude Opus 4.6", "opus-4.6"},
	{"anthropic", "claude-sonnet-4-6", "Claude Sonnet 4.6", "sonnet-4.6"},
	{"anthropic", "claude-haiku-4-5-20251001", "Claude Haiku 4.5", "haiku-4.5"},
	// OpenAI
	{"openai", "gpt-4o", "GPT-4o", "gpt-4o"},
	{"openai", "gpt-4o-mini", "GPT-4o Mini", "gpt-4o-mini"},
	{"openai", "o3-mini", "o3 Mini", "o3-mini"},
	// Gemini
	{"gemini", "gemini-2.0-flash", "Gemini 2.0 Flash", "gemini-2.0f"},
	{"gemini", "gemini-1.5-pro", "Gemini 1.5 Pro", "gemini-1.5p"},
	{"gemini", "gemini-1.5-flash", "Gemini 1.5 Flash", "gemini-1.5f"},
}

// ByID returns the model with the given ID, or zero value if not found.
func ByID(id string) (Model, bool) {
	for _, m := range All {
		if m.ID == id {
			return m, true
		}
	}
	return Model{}, false
}

// ShortName returns a short display name for a model ID.
func ShortName(id string) string {
	if m, ok := ByID(id); ok {
		return m.Short
	}
	if len(id) > 12 {
		return id[:12]
	}
	return id
}

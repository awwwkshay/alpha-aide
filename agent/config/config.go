package config

import (
	"bufio"
	"flag"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Provider  string
	Model     string
	BaseURL   string
	MaxTokens int64
	NoColor   bool

	AnthropicAPIKey string
	OpenAIAPIKey    string
	GeminiAPIKey    string
}

func Load() *Config {
	loadDotEnv()
	cfg := &Config{}

	flag.StringVar(&cfg.Provider, "provider", getEnv("AGENT_PROVIDER", "anthropic"), "LLM provider: anthropic, openai, or gemini")
	flag.StringVar(&cfg.Model, "model", getEnv("AGENT_MODEL", "claude-haiku-4-5-20251001"), "Model name")
	flag.StringVar(&cfg.BaseURL, "baseurl", getEnv("OPENAI_BASE_URL", ""), "Base URL for OpenAI-compatible providers (leave empty for api.openai.com)")
	maxTokens := flag.Int64("max-tokens", 8096, "Maximum tokens in response")
	flag.BoolVar(&cfg.NoColor, "no-color", os.Getenv("NO_COLOR") != "", "Disable ANSI colors")
	flag.Parse()

	cfg.MaxTokens = *maxTokens
	cfg.AnthropicAPIKey = os.Getenv("ANTHROPIC_API_KEY")
	cfg.OpenAIAPIKey = os.Getenv("OPENAI_API_KEY")
	cfg.GeminiAPIKey = os.Getenv("GEMINI_API_KEY")

	return cfg
}

func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// loadDotEnv loads env vars from ~/.alpha-aide/.env (global) and the nearest
// project .env (cwd walk-up). Project values override global; shell env wins over both.
func loadDotEnv() {
	globalVars := readEnvFile(findGlobalDotEnv())
	projectVars := readEnvFile(findProjectDotEnv())

	// merge: project takes precedence over global
	merged := globalVars
	for k, v := range projectVars {
		merged[k] = v
	}

	// only set keys not already exported in shell
	for k, v := range merged {
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}

// findGlobalDotEnv returns the path to ~/.alpha-aide/.env if it exists.
func findGlobalDotEnv() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	p := filepath.Join(home, ".alpha-aide", ".env")
	if _, err := os.Stat(p); err == nil {
		return p
	}
	return ""
}

// findProjectDotEnv walks up from the working directory looking for a .env file.
func findProjectDotEnv() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	for {
		candidate := filepath.Join(dir, ".env")
		if _, err := os.Stat(candidate); err == nil {
			return candidate
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// readEnvFile parses a .env file into a map. Returns an empty map if path is empty or unreadable.
func readEnvFile(path string) map[string]string {
	vars := make(map[string]string)
	if path == "" {
		return vars
	}
	f, err := os.Open(path)
	if err != nil {
		return vars
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if key != "" {
			vars[key] = value
		}
	}
	return vars
}

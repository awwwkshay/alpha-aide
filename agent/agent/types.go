package agent

import "github.com/awwwkshay/alpha-aide/llm"

// Re-export llm types so callers don't need a direct dependency on the llm package.
type (
	StreamChunk = llm.StreamChunk
	Tool        = llm.Tool
	Provider    = llm.Provider
)

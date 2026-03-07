package gemini

import (
	"context"
	"fmt"

	"github.com/awwwkshay/alpha-aide/llm"
	"google.golang.org/genai"
)

type Provider struct {
	apiKey    string
	model     string
	maxTokens int32
}

func New(apiKey, model string, maxTokens int64) *Provider {
	return &Provider{apiKey: apiKey, model: model, maxTokens: int32(maxTokens)}
}

func (p *Provider) Complete(ctx context.Context, systemPrompt string, messages []llm.Message, tools []llm.Tool) (<-chan llm.StreamChunk, error) {
	if len(messages) == 0 {
		return nil, fmt.Errorf("no messages provided")
	}

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  p.apiKey,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return nil, fmt.Errorf("create gemini client: %w", err)
	}

	config := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{{Text: systemPrompt}},
		},
		MaxOutputTokens: p.maxTokens,
	}
	if len(tools) > 0 {
		config.Tools = convertTools(tools)
	}

	contents := convertMessages(messages)

	ch := make(chan llm.StreamChunk, 64)
	go func() {
		defer close(ch)

		for resp, err := range client.Models.GenerateContentStream(ctx, p.model, contents, config) {
			if err != nil {
				ch <- llm.StreamChunk{Type: "error", Err: err}
				return
			}
			for _, candidate := range resp.Candidates {
				if candidate.Content == nil {
					continue
				}
				for i, part := range candidate.Content.Parts {
					if part.Text != "" {
						ch <- llm.StreamChunk{Type: "text_delta", Text: part.Text}
					}
					if part.FunctionCall != nil {
						fc := part.FunctionCall
						id := fc.ID
						if id == "" {
							id = fmt.Sprintf("%s_%d", fc.Name, i)
						}
						ch <- llm.StreamChunk{
							Type: "tool_use_start",
							ToolCall: &llm.ToolCall{
								ID:    id,
								Name:  fc.Name,
								Input: fc.Args,
							},
						}
						ch <- llm.StreamChunk{Type: "tool_use_end"}
					}
				}
			}
		}

		ch <- llm.StreamChunk{Type: "done"}
	}()

	return ch, nil
}

func convertMessages(messages []llm.Message) []*genai.Content {
	result := make([]*genai.Content, 0, len(messages))
	for _, m := range messages {
		content := &genai.Content{}
		switch m.Role {
		case llm.RoleUser:
			content.Role = "user"
		case llm.RoleAssistant:
			content.Role = "model"
		}
		for _, c := range m.Content {
			switch c.Type {
			case "text":
				content.Parts = append(content.Parts, genai.NewPartFromText(c.Text))
			case "tool_use":
				content.Parts = append(content.Parts, genai.NewPartFromFunctionCall(c.ToolName, c.ToolInput))
			case "tool_result":
				resp := map[string]any{"result": c.Content}
				if c.IsError {
					resp = map[string]any{"error": c.Content}
				}
				content.Parts = append(content.Parts, genai.NewPartFromFunctionResponse(c.ToolName, resp))
			}
		}
		if len(content.Parts) > 0 {
			result = append(result, content)
		}
	}
	return result
}

func convertTools(tools []llm.Tool) []*genai.Tool {
	decls := make([]*genai.FunctionDeclaration, 0, len(tools))
	for _, t := range tools {
		decls = append(decls, &genai.FunctionDeclaration{
			Name:                 t.Name(),
			Description:          t.Description(),
			ParametersJsonSchema: t.InputSchema(),
		})
	}
	return []*genai.Tool{{FunctionDeclarations: decls}}
}

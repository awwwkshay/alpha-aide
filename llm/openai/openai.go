package openai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/awwwkshay/alpha-aide/llm"
	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

type Provider struct {
	client    *openai.Client
	model     string
	maxTokens int64
}

func New(apiKey, baseURL, model string, maxTokens int64) *Provider {
	opts := []option.RequestOption{option.WithAPIKey(apiKey)}
	if baseURL != "" {
		opts = append(opts, option.WithBaseURL(baseURL))
	}
	client := openai.NewClient(opts...)
	return &Provider{client: client, model: model, maxTokens: maxTokens}
}

func (p *Provider) Complete(ctx context.Context, systemPrompt string, messages []llm.Message, tools []llm.Tool) (<-chan llm.StreamChunk, error) {
	apiMessages := convertMessages(systemPrompt, messages)
	apiTools := convertTools(tools)

	params := openai.ChatCompletionNewParams{
		Model:     openai.F(p.model),
		MaxTokens: openai.F(p.maxTokens),
		Messages:  openai.F(apiMessages),
	}
	if len(apiTools) > 0 {
		params.Tools = openai.F(apiTools)
	}

	stream := p.client.Chat.Completions.NewStreaming(ctx, params)

	ch := make(chan llm.StreamChunk, 64)
	go func() {
		defer close(ch)
		defer stream.Close()

		// Track tool calls accumulating across chunks
		type accumTC struct {
			id    string
			name  string
			input string
		}
		toolCalls := map[int]*accumTC{}

		for stream.Next() {
			chunk := stream.Current()
			for _, choice := range chunk.Choices {
				delta := choice.Delta
				if delta.Content != "" {
					ch <- llm.StreamChunk{Type: "text_delta", Text: delta.Content}
				}
				for _, tc := range delta.ToolCalls {
					idx := int(tc.Index)
					if _, ok := toolCalls[idx]; !ok {
						toolCalls[idx] = &accumTC{}
					}
					a := toolCalls[idx]
					if tc.ID != "" && a.id == "" {
						a.id = fmt.Sprintf("call_%d", idx)
					}
					if tc.Function.Name != "" {
						a.name = tc.Function.Name
					}
					a.input += tc.Function.Arguments
				}
			}
		}

		if err := stream.Err(); err != nil {
			ch <- llm.StreamChunk{Type: "error", Err: err}
			return
		}

		// Emit tool calls in index order
		for i := 0; i < len(toolCalls); i++ {
			a, ok := toolCalls[i]
			if !ok {
				break
			}
			var input map[string]any
			if a.input != "" {
				if err := json.Unmarshal([]byte(a.input), &input); err != nil {
					ch <- llm.StreamChunk{Type: "error", Err: fmt.Errorf("parse tool input: %w", err)}
					return
				}
			}
			ch <- llm.StreamChunk{
				Type: "tool_use_start",
				ToolCall: &llm.ToolCall{
					ID:    a.id,
					Name:  a.name,
					Input: input,
				},
			}
			ch <- llm.StreamChunk{Type: "tool_use_end"}
		}

		ch <- llm.StreamChunk{Type: "done"}
	}()

	return ch, nil
}

func convertMessages(systemPrompt string, messages []llm.Message) []openai.ChatCompletionMessageParamUnion {
	result := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages)+1)
	result = append(result, openai.SystemMessage(systemPrompt))

	for _, m := range messages {
		switch m.Role {
		case llm.RoleUser:
			// Check if this is a tool result message
			var toolResults []openai.ChatCompletionToolMessageParam
			var textParts []string
			for _, c := range m.Content {
				switch c.Type {
				case "tool_result":
					toolResults = append(toolResults, openai.ToolMessage(c.ToolUseID, c.Content))
				case "text":
					textParts = append(textParts, c.Text)
				}
			}
			for _, tr := range toolResults {
				result = append(result, tr)
			}
			if len(textParts) > 0 {
				for _, t := range textParts {
					result = append(result, openai.UserMessage(t))
				}
			}
		case llm.RoleAssistant:
			var text string
			var toolCalls []openai.ChatCompletionMessageToolCallParam
			for _, c := range m.Content {
				switch c.Type {
				case "text":
					text += c.Text
				case "tool_use":
					inputJSON, _ := json.Marshal(c.ToolInput)
					toolCalls = append(toolCalls, openai.ChatCompletionMessageToolCallParam{
						ID:   openai.F(c.ToolUseID),
						Type: openai.F(openai.ChatCompletionMessageToolCallTypeFunction),
						Function: openai.F(openai.ChatCompletionMessageToolCallFunctionParam{
							Name:      openai.F(c.ToolName),
							Arguments: openai.F(string(inputJSON)),
						}),
					})
				}
			}
			msg := openai.ChatCompletionAssistantMessageParam{
				Role: openai.F(openai.ChatCompletionAssistantMessageParamRoleAssistant),
			}
			if text != "" {
				msg.Content = openai.F([]openai.ChatCompletionAssistantMessageParamContentUnion{
					openai.TextPart(text),
				})
			}
			if len(toolCalls) > 0 {
				msg.ToolCalls = openai.F(toolCalls)
			}
			result = append(result, msg)
		}
	}
	return result
}


func convertTools(tools []llm.Tool) []openai.ChatCompletionToolParam {
	result := make([]openai.ChatCompletionToolParam, 0, len(tools))
	for _, t := range tools {
		schema := t.InputSchema()
		result = append(result, openai.ChatCompletionToolParam{
			Type: openai.F(openai.ChatCompletionToolTypeFunction),
			Function: openai.F(openai.FunctionDefinitionParam{
				Name:        openai.F(t.Name()),
				Description: openai.F(t.Description()),
				Parameters:  openai.F(openai.FunctionParameters(schema)),
			}),
		})
	}
	return result
}

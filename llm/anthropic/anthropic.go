package anthropic

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/awwwkshay/alpha-aide/llm"
	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

type Provider struct {
	client    *anthropic.Client
	model     string
	maxTokens int64
}

func New(apiKey, model string, maxTokens int64) *Provider {
	client := anthropic.NewClient(option.WithAPIKey(apiKey))
	return &Provider{client: client, model: model, maxTokens: maxTokens}
}

func (p *Provider) Complete(ctx context.Context, systemPrompt string, messages []llm.Message, tools []llm.Tool) (<-chan llm.StreamChunk, error) {
	apiMessages, err := convertMessages(messages)
	if err != nil {
		return nil, err
	}

	params := anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(p.model)),
		MaxTokens: anthropic.F(p.maxTokens),
		System: anthropic.F([]anthropic.TextBlockParam{
			anthropic.NewTextBlock(systemPrompt),
		}),
		Messages: anthropic.F(apiMessages),
	}

	if len(tools) > 0 {
		apiTools := convertTools(tools)
		params.Tools = anthropic.F(apiTools)
	}

	stream := p.client.Messages.NewStreaming(ctx, params)

	ch := make(chan llm.StreamChunk, 64)
	go func() {
		defer close(ch)
		defer stream.Close()

		var currentToolCall *llm.ToolCall
		var inputAccum string

		for stream.Next() {
			event := stream.Current()
			switch e := event.AsUnion().(type) {
			case anthropic.ContentBlockStartEvent:
				block := e.ContentBlock
				if block.Type == "tool_use" {
					currentToolCall = &llm.ToolCall{
						ID:   block.ID,
						Name: block.Name,
					}
					inputAccum = ""
				}
			case anthropic.ContentBlockDeltaEvent:
				switch delta := e.Delta.AsUnion().(type) {
				case anthropic.TextDelta:
					ch <- llm.StreamChunk{Type: "text_delta", Text: delta.Text}
				case anthropic.InputJSONDelta:
					inputAccum += delta.PartialJSON
				}
			case anthropic.ContentBlockStopEvent:
				if currentToolCall != nil {
					var input map[string]any
					if inputAccum != "" {
						if err := json.Unmarshal([]byte(inputAccum), &input); err != nil {
							ch <- llm.StreamChunk{Type: "error", Err: fmt.Errorf("parse tool input: %w", err)}
							return
						}
					}
					currentToolCall.Input = input
					ch <- llm.StreamChunk{Type: "tool_use_start", ToolCall: currentToolCall}
					ch <- llm.StreamChunk{Type: "tool_use_end"}
					currentToolCall = nil
					inputAccum = ""
				}
			}
		}

		if err := stream.Err(); err != nil {
			ch <- llm.StreamChunk{Type: "error", Err: err}
			return
		}

		ch <- llm.StreamChunk{Type: "done"}
	}()

	return ch, nil
}

func convertMessages(messages []llm.Message) ([]anthropic.MessageParam, error) {
	result := make([]anthropic.MessageParam, 0, len(messages))
	for _, m := range messages {
		blocks := make([]anthropic.ContentBlockParamUnion, 0, len(m.Content))
		for _, c := range m.Content {
			switch c.Type {
			case "text":
				blocks = append(blocks, anthropic.NewTextBlock(c.Text))
			case "tool_use":
				blocks = append(blocks, anthropic.NewToolUseBlockParam(c.ToolUseID, c.ToolName, c.ToolInput))
			case "tool_result":
				blocks = append(blocks, anthropic.NewToolResultBlock(c.ToolUseID, c.Content, c.IsError))
			}
		}
		switch m.Role {
		case llm.RoleUser:
			result = append(result, anthropic.NewUserMessage(blocks...))
		case llm.RoleAssistant:
			result = append(result, anthropic.NewAssistantMessage(blocks...))
		}
	}
	return result, nil
}

func convertTools(tools []llm.Tool) []anthropic.ToolUnionUnionParam {
	result := make([]anthropic.ToolUnionUnionParam, 0, len(tools))
	for _, t := range tools {
		schema := t.InputSchema()
		result = append(result, anthropic.ToolUnionParam{
			Name:        anthropic.F(t.Name()),
			Description: anthropic.F(t.Description()),
			InputSchema: anthropic.F[interface{}](schema),
		})
	}
	return result
}

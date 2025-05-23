package functioncall

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/flc1125/ai-agent-share/infra"
	"github.com/sashabaranov/go-openai"
	"github.com/sashabaranov/go-openai/jsonschema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTool(t *testing.T) openai.Tool {
	return openai.Tool{
		Type: openai.ToolTypeFunction,
		Function: &openai.FunctionDefinition{
			Name:        "get_current_time",
			Description: "Get current system time.",
			Parameters: jsonschema.Definition{
				Type: jsonschema.Object,
				Properties: map[string]jsonschema.Definition{
					"format": {
						Type:        jsonschema.String,
						Description: "The format of the time. The format configuration for time is the same as that in the Go language.",
					},
				},
				Required: []string{"format"},
			},
		},
	}
}

type GetCurrentTimeArgs struct {
	Format string `json:"format"`
}

func callback(t *testing.T, functionName string, functionArgs string) string {
	switch functionName {
	case "get_current_time":
		var args GetCurrentTimeArgs
		require.NoError(t, json.Unmarshal([]byte(functionArgs), &args))
		if args.Format == "" {
			args.Format = time.DateTime
		}
		return fmt.Sprintf("现在的时间是：%s", time.Now().Format(args.Format))
	default:
		t.Logf("Unknown function: %s", functionName)
		return fmt.Sprintf("Unknown function: %s", functionName)
	}
}

func TestFunctionCall(t *testing.T) {
	client := infra.NewArtModelBaseOpenAIProtocol(t)

	var invoker = func(ctx context.Context, messages []openai.ChatCompletionMessage) (openai.ChatCompletionResponse, error) {
		assert.NotEmpty(t, messages)
		t.Logf("message length: %d", len(messages))

		return client.CreateChatCompletion(
			ctx,
			openai.ChatCompletionRequest{
				Model:    infra.DefaultModel,
				Messages: messages,
				Tools: []openai.Tool{
					newTool(t),
				},
				Temperature: 0.1,
			},
		)
	}

	messages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "现在几点了？",
		},
	}

	for {
		resp, err := invoker(t.Context(), messages)
		require.NoError(t, err)
		require.NotEmpty(t, resp.Choices)
		t.Logf("reply: %s", resp.Choices[0].Message.Content)
		// spew.Dump(resp)

		require.NotEmpty(t, resp.Choices)

		// If the model has finished its response, break the loop
		if resp.Choices[0].FinishReason == openai.FinishReasonStop {
			break
		}

		// if call tools
		if resp.Choices[0].FinishReason == openai.FinishReasonToolCalls {
			assert.NotEmpty(t, resp.Choices[0].Message.ToolCalls)

			// append message
			messages = append(messages, resp.Choices[0].Message)

			for _, toolCall := range resp.Choices[0].Message.ToolCalls {
				result := callback(t, toolCall.Function.Name, toolCall.Function.Arguments)
				t.Logf("Tool call: %s, args: %s, result: %s",
					toolCall.Function.Name, toolCall.Function.Arguments, result)

				// append message
				messages = append(messages, openai.ChatCompletionMessage{
					Role:       openai.ChatMessageRoleTool,
					ToolCallID: toolCall.ID,
					Content:    result,
				})
			}

			continue
		}
	}
}

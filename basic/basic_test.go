package basic

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/schema"
	"github.com/flc1125/ai-agent-share/infra"
	"github.com/flc1125/ai-agent-share/util"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestChatModel_NewArkModel(t *testing.T) {
	cm := infra.NewArkModel(t)

	reader, err := cm.Stream(t.Context(), []*schema.Message{
		schema.SystemMessage(`你是一个 Ping/Pong 服务器。
当我给你发送 Ping 时，你只需要回复 Pong 给我。
不要过多的理解，也不要过多的回复无关紧要的信息。`),
		schema.UserMessage("Ping"),
	})
	require.NoError(t, err)
	defer reader.Close()

	util.PrintContentByReader(t, reader)
}

func TestChatModel_NewArtModelBaseOpenAIProtocol(t *testing.T) {
	client := infra.NewArtModelBaseOpenAIProtocol(t)

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: infra.DefaultModel,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: "你是谁？",
				},
			},
		},
	)
	require.NoError(t, err)
	assert.NotEmpty(t, resp.Choices)
	assert.NotEmpty(t, resp.Choices[0].Message.Content)
	t.Logf("reply: %s", resp.Choices[0].Message.Content)
	// spew.Dump(resp)
}

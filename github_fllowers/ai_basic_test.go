package github_fllowers

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/schema"
	"github.com/flc1125/ai-agent-share/infra"
	"github.com/flc1125/ai-agent-share/tools/getrequester"
	"github.com/flc1125/ai-agent-share/tools/githubfollower"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAI_Basic(t *testing.T) {
	t.Logf("reply: %s", getFlowersByBasic(t, "获取 https://flc.io/eino-deepseek-rss/ 的作者的Github账号所对应的粉丝数量"))
}

func getFlowersByBasic(t *testing.T, prompt string) string {
	cm := infra.NewArkModel(t)

	// Initialize the tools
	getRequester, err := getrequester.NewTool()
	require.NoError(t, err)

	githubFollower, err := githubfollower.NewTool()
	require.NoError(t, err)

	tools := []tool.InvokableTool{getRequester, githubFollower}
	toolInfos := make([]*schema.ToolInfo, 0, len(tools))
	for _, it := range tools {
		info, err := it.Info(t.Context())
		require.NoError(t, err)
		toolInfos = append(toolInfos, info)
	}

	// Set the tools in the model
	cm, err = cm.WithTools(toolInfos)
	require.NoError(t, err)

	// define tool callback
	toolCallback := func(ctx context.Context, toolName string, args string) (string, error) {
		for _, it := range tools {
			info, err := it.Info(ctx)
			require.NoError(t, err)
			if info.Name != toolName {
				continue
			}

			// call the tool
			return it.InvokableRun(ctx, args)
		}
		return "", nil
	}

	// define invoker
	round := 1
	invoker := func(ctx context.Context, messages []*schema.Message) (*schema.Message, error) {
		assert.NotEmpty(t, messages)

		var buf bytes.Buffer
		for _, message := range messages {
			buf.WriteString(fmt.Sprintf("%s: %s\n", message.Role, message.Content))
		}
		t.Logf("第 %d 次调用: \n%s\n\n", round, buf.String())
		round++

		return cm.Generate(ctx, messages)
	}

	// define messages
	messages := []*schema.Message{
		schema.SystemMessage(prompt),
	}

	for {
		// call the model
		message, err := invoker(t.Context(), messages)
		require.NoError(t, err)
		require.NotNil(t, message)

		// append message
		messages = append(messages, message)

		if len(message.ToolCalls) > 0 {
			for _, toolCall := range message.ToolCalls {
				result, err := toolCallback(t.Context(), toolCall.Function.Name, toolCall.Function.Arguments)
				require.NoError(t, err)
				// t.Logf("==============")
				// t.Logf("Tool call: %s, args: %s, result: %s",
				// 	toolCall.Function.Name, toolCall.Function.Arguments, result)
				// t.Logf("==============")

				// append message
				messages = append(messages, schema.ToolMessage(result, toolCall.ID))
			}
		} else {
			return message.Content
		}
	}
}

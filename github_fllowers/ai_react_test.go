package github_fllowers

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/flc1125/ai-agent-share/infra"
	"github.com/flc1125/ai-agent-share/tools/getrequester"
	"github.com/flc1125/ai-agent-share/tools/githubfollower"
	"github.com/stretchr/testify/require"
)

func TestAI_ReAct(t *testing.T) {
	t.Logf("reply: %s", getFlowersByModel(t, "获取 https://flc.io/eino-deepseek-rss/ 的作者的Github账号所对应的粉丝数量"))
}

func getFlowersByModel(t *testing.T, prompt string) string {
	cm := infra.NewQwenModel(t)

	getRequester, err := getrequester.NewTool()
	require.NoError(t, err)

	githubFollower, err := githubfollower.NewTool()
	require.NoError(t, err)

	agent, err := react.NewAgent(t.Context(), &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				getRequester,
				githubFollower,
			},
		},
		MaxStep: 10,
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			// spew.Dump(input, "==========")
			return input
		},
	})
	require.NoError(t, err)

	message, err := agent.Generate(
		t.Context(),
		[]*schema.Message{
			schema.SystemMessage(prompt),
		},
	)
	require.NoError(t, err)
	require.NotNil(t, message)
	return message.Content
}

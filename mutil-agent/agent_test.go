package mutil_agent

import (
	"context"
	"testing"

	"github.com/cloudwego/eino/flow/agent"
	"github.com/cloudwego/eino/flow/agent/multiagent/host"
	"github.com/cloudwego/eino/schema"
	"github.com/davecgh/go-spew/spew"
	"github.com/flc1125/ai-agent-share/infra"
	"github.com/stretchr/testify/require"
)

func TestMutilAgent(t *testing.T) {
	cm := infra.NewQwenModel(t)

	multiAgent, err := host.NewMultiAgent(t.Context(), &host.MultiAgentConfig{
		Host: host.Host{
			ToolCallingModel: cm,
			SystemPrompt:     "You are a multi-agent system, and you can call other agents to complete the task.",
		},
		Specialists: []*host.Specialist{
			{
				AgentMeta: host.AgentMeta{
					Name:        "get_requester",
					IntendedUse: "Get the request information of the URL",
				},
				ChatModel: cm,
				Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.Message, err error) {
					return schema.UserMessage("Github: flc1125"), nil
				},
			},
			{
				AgentMeta: host.AgentMeta{
					Name:        "github_follower",
					IntendedUse: "Get the followers count of the github user",
				},
				ChatModel: cm,
				Invokable: func(ctx context.Context, input []*schema.Message, opts ...agent.AgentOption) (output *schema.Message, err error) {
					return schema.UserMessage("followers: 100"), nil
				},
			},
		},
	})
	require.NoError(t, err)

	message, err := multiAgent.Generate(t.Context(), []*schema.Message{
		schema.SystemMessage("获取 https://flc.io/eino-deepseek-rss/ 的作者的 Github 账号所对应的粉丝数量"),
	})
	require.NoError(t, err)
	spew.Dump(message)
}

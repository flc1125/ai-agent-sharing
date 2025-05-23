package micro_agent

import (
	"bytes"
	"context"
	"fmt"
	"testing"

	"github.com/cloudwego/eino/components/tool"
	"github.com/cloudwego/eino/components/tool/utils"
	"github.com/cloudwego/eino/compose"
	"github.com/cloudwego/eino/flow/agent/react"
	"github.com/cloudwego/eino/schema"
	"github.com/flc1125/ai-agent-share/infra"
	"github.com/stretchr/testify/require"
)

type getEmployeeNameInput struct {
	EmployeeID string `json:"employee_id" jsonschema_description:"The ID of the employee."`
}

type getEmployeeNameOutput struct {
	EmployeeName string `json:"employee_name" jsonschema_description:"The name of the employee."`
}

func newTool[I, O any](t *testing.T, name, desc string, output O) tool.InvokableTool {
	it, err := utils.InferTool(name, desc, func(ctx context.Context, input I) (O, error) {
		return output, nil
	})
	require.NoError(t, err)
	return it
}

func TestMicroAgent(t *testing.T) {
	cm := infra.NewArkModel(t)

	reactAgent, err := react.NewAgent(t.Context(), &react.AgentConfig{
		ToolCallingModel: cm,
		ToolsConfig: compose.ToolsNodeConfig{
			Tools: []tool.BaseTool{
				newTool[getEmployeeNameInput, getEmployeeNameOutput](t,
					"Get user name",
					"Get the user name",
					getEmployeeNameOutput{EmployeeName: "李四"},
				),
				newTool[getEmployeeNameInput, getEmployeeNameOutput](t,
					"Get employee name",
					"Get the employee name",
					getEmployeeNameOutput{EmployeeName: "张三"},
				),
			},
		},
		MessageModifier: func(ctx context.Context, input []*schema.Message) []*schema.Message {
			var buf bytes.Buffer
			for _, msg := range input {
				buf.WriteString(fmt.Sprintf("%s: %s\n", msg.Role, msg.Content))
			}
			t.Logf("input: %s", buf.String())

			return input
		},
	})
	require.NoError(t, err)

	message, err := reactAgent.Generate(t.Context(), []*schema.Message{
		schema.UserMessage("帮我看看123的名字"),
	})
	require.NoError(t, err)
	require.NotNil(t, message)
	t.Logf("reply: %s", message.Content)
}

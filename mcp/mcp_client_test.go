package mcp

import (
	"context"
	"fmt"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/mark3labs/mcp-go/client"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMCPClient(t *testing.T) {
	// https://github.com/modelcontextprotocol/servers/tree/main/src/redis
	c, err := client.NewStdioMCPClient(
		"npx",
		[]string{}, // Empty ENV
		"-y",
		"@modelcontextprotocol/server-redis",
		"redis://localhost:6379",
	)
	require.NoError(t, err)
	defer c.Close()

	// Initialize the client
	fmt.Println("Initializing client...")
	initRequest := mcp.InitializeRequest{}
	initRequest.Params.ProtocolVersion = mcp.LATEST_PROTOCOL_VERSION
	initRequest.Params.ClientInfo = mcp.Implementation{
		Name: "example-client",
	}

	initResult, err := c.Initialize(t.Context(), initRequest)
	require.NoError(t, err)
	spew.Dump(initResult)
	println("\n\n==========INIT END============\n\n")

	// list tools
	tools, err := c.ListTools(t.Context(), mcp.ListToolsRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, tools)
	spew.Dump(tools)
	println("\n\n==========LIST TOOLS END============\n\n")

	// call set
	setCallToolRequest := mcp.CallToolRequest{}
	setCallToolRequest.Params.Name = "set"
	setCallToolRequest.Params.Arguments = map[string]any{
		"key":   "foo",
		"value": "bar",
	}
	setCallToolResult, err := c.CallTool(t.Context(), setCallToolRequest)
	require.NoError(t, err)
	spew.Dump(setCallToolResult)
	println("\n\n==========SET END============\n\n")

	// call get
	getCallToolRequest := mcp.CallToolRequest{}
	getCallToolRequest.Params.Name = "get"
	getCallToolRequest.Params.Arguments = map[string]any{
		"key": "foo",
	}
	getCallToolResult, err := c.CallTool(t.Context(), getCallToolRequest)
	require.NoError(t, err)
	spew.Dump(getCallToolResult)
	println("\n\n==========GET END============\n\n")
}

func newMCPServer(t *testing.T) *server.MCPServer {
	srv := server.NewMCPServer(
		"example-server",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	helloWorldTool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)

	srv.AddTool(helloWorldTool, func(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		t.Logf("Tool call, the args: %v", request.GetArguments())

		name, ok := request.GetArguments()["name"].(string)
		if !ok {
			return mcp.NewToolResultError("name must be a string"), nil
		}

		return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
	})

	return srv
}

func TestMCPServer_SSE(t *testing.T) {
	srv := newMCPServer(t)

	assert.NoError(t, server.NewSSEServer(srv).Start(":8200"))
}

func TestMCPServer_Stream(t *testing.T) {
	srv := newMCPServer(t)

	assert.NoError(t, server.NewStreamableHTTPServer(srv).Start(":8300"))
}

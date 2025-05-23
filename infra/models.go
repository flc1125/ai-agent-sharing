package infra

import (
	"cmp"
	"os"
	"testing"

	"github.com/cloudwego/eino-ext/components/model/ark"
	"github.com/cloudwego/eino-ext/components/model/qwen"
	"github.com/cloudwego/eino/components/model"
	"github.com/sashabaranov/go-openai"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	arkAPIKey  = cmp.Or(os.Getenv("ARK_API_KEY"), "key")
	qwenAPIKey = cmp.Or(os.Getenv("QWEN_API_KEY"), "key")
	// DefaultModel = "doubao-pro-256k-241115"
	DefaultModel = "doubao-1-5-lite-32k-250115"
)

func NewArkModel(tb testing.TB) model.ToolCallingChatModel {
	cm, err := ark.NewChatModel(tb.Context(), &ark.ChatModelConfig{
		APIKey: arkAPIKey,
		Model:  DefaultModel,
	})
	require.NoError(tb, err)
	assert.NotNil(tb, cm)
	return cm
}

func NewArtModelBaseOpenAIProtocol(tb testing.TB) *openai.Client {
	cfg := openai.DefaultConfig(arkAPIKey)
	cfg.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"
	return openai.NewClientWithConfig(cfg)
}

func NewQwenModel(tb testing.TB) model.ToolCallingChatModel {
	cm, err := qwen.NewChatModel(tb.Context(), &qwen.ChatModelConfig{
		APIKey:  qwenAPIKey,
		Model:   "qwen-plus-latest",
		BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1",
	})
	require.NoError(tb, err)
	assert.NotNil(tb, cm)
	return cm
}

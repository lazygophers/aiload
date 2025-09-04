package channel

import (
	"errors"
	"time"

	"github.com/lazygophers/aiload"
)

type GetModelListRsp_Data struct {
	Id          string             `json:"id,omitempty"`
	Name        string             `json:"name,omitempty"`
	Object      string             `json:"object,omitempty"`
	Created     int                `json:"created,omitempty"`
	OwnedBy     string             `json:"owned_by,omitempty"`
	Model       string             `json:"model,omitempty"`
	ModifiedAt  time.Time          `json:"modified_at,omitempty"`
	Size        int64              `json:"size,omitempty"`
	Digest      string             `json:"digest,omitempty"`
	Details     OllamaModelDetails `json:"details,omitempty"`
	BaseModelID string             `json:"baseModelId,omitempty"`
	// Version 模型的版本号。
	Version string `json:"version,omitempty"`
	// DisplayName 模型在 UI 中展示的名称。
	DisplayName string `json:"displayName,omitempty"`
	// Description 模型的详细描述。
	Description string `json:"description,omitempty"`
	// InputTokenLimit 模型支持的最大输入 token 数。
	InputTokenLimit int `json:"inputTokenLimit,omitempty"`
	// OutputTokenLimit 模型支持的最大输出 token 数。
	OutputTokenLimit int `json:"outputTokenLimit,omitempty"`
	// SupportedGenerationMethods 支持的生成方法列表。
	SupportedGenerationMethods []string `json:"supportedGenerationMethods,omitempty"`
	// Thinking 指示模型是否处于“思考”状态。
	Thinking bool `json:"thinking,omitempty"`
	// Temperature 模型的温度参数。
	Temperature int `json:"temperature,omitempty"`
	// MaxTemperature 模型的最大温度参数。
	MaxTemperature int `json:"maxTemperature,omitempty"`
	// TopP 模型的 Top-p 采样参数。
	TopP int `json:"topP,omitempty"`
	// TopK 模型的 Top-k 采样参数。
	TopK int `json:"topK,omitempty"`
}

type GetModelListRsp struct {
	Object        string                  `json:"object,omitempty"`
	Data          []*GetModelListRsp_Data `json:"data,omitempty"`
	NextPageToken string                  `json:"nextPageToken,omitempty"`
}

type ChatMessage struct {
	Role    string `json:"role,omitempty"`
	Content string `json:"content,omitempty"`
}

type ChatCompletionReq struct {
	Model    string         `json:"model,omitempty"`
	Stream   bool           `json:"stream,omitempty"`
	Messages []*ChatMessage `json:"messages,omitempty"`
}

type ChatCompletionRspItem struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index int `json:"index"`
		Delta struct {
			Content          string      `json:"content"`
			ReasoningContent interface{} `json:"reasoning_content"`
			Role             string      `json:"role"`
		} `json:"delta"`
		FinishReason interface{} `json:"finish_reason"`
	} `json:"choices"`
	SystemFingerprint string `json:"system_fingerprint"`
	Usage             struct {
		PromptTokens            int `json:"prompt_tokens"`
		CompletionTokens        int `json:"completion_tokens"`
		TotalTokens             int `json:"total_tokens"`
		CompletionTokensDetails struct {
			ReasoningTokens int `json:"reasoning_tokens"`
		} `json:"completion_tokens_details"`
	} `json:"usage"`
}

type Channeler interface {
	GetModelList() (*GetModelListRsp, error)
	ChatCompletionsAsync(req *ChatCompletionReq, callback func(item *ChatCompletionRspItem) error) error
}

func NewChannel(channel *aiload.ModelChannel) (Channeler, error) {
	switch channel.Platform {
	case aiload.Platform_Siliconflow:
		return NewSiliconFlowAdapter(channel), nil
	case aiload.Platform_Gemini:
		return NewGeminiAdapter(channel), nil
	case aiload.Platform_Ollama:
		return NewOllamaAdapter(channel), nil
	case aiload.Platform_OpenAi:
		return NewOpenAiAdapter(channel), nil
	case aiload.Platform_OpenAiCompatible:
		return NewOpenAiCompatibleAdapter(channel), nil
	default:
		return nil, errors.New("unsupported platform")
	}
}

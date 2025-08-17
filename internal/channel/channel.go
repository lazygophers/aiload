package channel

import (
	"context"
	"errors"
	"time"

	"github.com/lazygophers/aiload"
)

type GetModelListReq struct {
	Type      string
	SubType   string
	PageSize  int
	PageToken string
}

type GetModelListRsp struct {
	Object string `json:"object"`
	Data   []struct {
		Id          string             `json:"id"`
		Name        string             `json:"name"`
		Object      string             `json:"object"`
		Created     int                `json:"created"`
		OwnedBy     string             `json:"owned_by"`
		Model       string             `json:"model"`
		ModifiedAt  time.Time          `json:"modified_at"`
		Size        int64              `json:"size"`
		Digest      string             `json:"digest"`
		Details     OllamaModelDetails `json:"details"`
		BaseModelID string             `json:"baseModelId"`
		// Version 模型的版本号。
		Version string `json:"version"`
		// DisplayName 模型在 UI 中展示的名称。
		DisplayName string `json:"displayName"`
		// Description 模型的详细描述。
		Description string `json:"description"`
		// InputTokenLimit 模型支持的最大输入 token 数。
		InputTokenLimit int `json:"inputTokenLimit"`
		// OutputTokenLimit 模型支持的最大输出 token 数。
		OutputTokenLimit int `json:"outputTokenLimit"`
		// SupportedGenerationMethods 支持的生成方法列表。
		SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
		// Thinking 指示模型是否处于“思考”状态。
		Thinking bool `json:"thinking"`
		// Temperature 模型的温度参数。
		Temperature int `json:"temperature"`
		// MaxTemperature 模型的最大温度参数。
		MaxTemperature int `json:"maxTemperature"`
		// TopP 模型的 Top-p 采样参数。
		TopP int `json:"topP"`
		// TopK 模型的 Top-k 采样参数。
		TopK int `json:"topK"`
	} `json:"data"`
	NextPageToken string `json:"nextPageToken"`
}

type Channeler interface {
	GetModelList(context.Context, *GetModelListReq) (*GetModelListRsp, error)
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

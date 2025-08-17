package channel

import (
	"context"
	"time"

	"github.com/lazygophers/aiload"
	"github.com/lazygophers/utils/anyx"
)

// OllamaAdapter 实现了 Channeler 接口，用于适配 Ollama 客户端。
type OllamaAdapter struct {
	client *Ollama
}

// NewOllamaAdapter 创建一个新的 Ollama 适配器实例。
func NewOllamaAdapter(channel *aiload.ModelChannel) *OllamaAdapter {
	return &OllamaAdapter{
		client: NewOllama(channel),
	}
}

// GetModelList 调用 Ollama 客户端的 GetLocalModelList 方法，并将结果转换为统一的 GetModelListRsp 格式。
func (a *OllamaAdapter) GetModelList(ctx context.Context, req *GetModelListReq) (*GetModelListRsp, error) {
	ollamaRsp, err := a.client.GetLocalModelList()
	if err != nil {
		return nil, err
	}

	rsp := &GetModelListRsp{
		Object: "list",
	}
	for _, model := range ollamaRsp.Models {
		rsp.Data = append(rsp.Data, struct {
			Id                         string             `json:"id"`
			Name                       string             `json:"name"`
			Object                     string             `json:"object"`
			Created                    int                `json:"created"`
			OwnedBy                    string             `json:"owned_by"`
			Model                      string             `json:"model"`
			ModifiedAt                 time.Time          `json:"modified_at"`
			Size                       int64              `json:"size"`
			Digest                     string             `json:"digest"`
			Details                    OllamaModelDetails `json:"details"`
			BaseModelID                string             `json:"baseModelId"`
			Version                    string             `json:"version"`
			DisplayName                string             `json:"displayName"`
			Description                string             `json:"description"`
			InputTokenLimit            int                `json:"inputTokenLimit"`
			OutputTokenLimit           int                `json:"outputTokenLimit"`
			SupportedGenerationMethods []string           `json:"supportedGenerationMethods"`
			Thinking                   bool               `json:"thinking"`
			Temperature                int                `json:"temperature"`
			MaxTemperature             int                `json:"maxTemperature"`
			TopP                       int                `json:"topP"`
			TopK                       int                `json:"topK"`
		}{
			Id:         model.Name, // Ollama 中 name 即是其唯一标识
			Name:       model.Name,
			Object:     "model",
			Model:      model.Model,
			ModifiedAt: model.ModifiedAt,
			Size:       model.Size,
			Digest:     model.Digest,
			Details:    model.Details,
		})
	}

	return rsp, nil
}

// GeminiAdapter 实现了 Channeler 接口，用于适配 Gemini 客户端。
type GeminiAdapter struct {
	client *Gemini
}

// NewGeminiAdapter 创建一个新的 Gemini 适配器实例。
func NewGeminiAdapter(channel *aiload.ModelChannel) *GeminiAdapter {
	return &GeminiAdapter{
		client: NewGemini(channel),
	}
}

// GetModelList 调用 Gemini 客户端的 GetModelList 方法，并将结果转换为统一的 GetModelListRsp 格式。
func (a *GeminiAdapter) GetModelList(ctx context.Context, req *GetModelListReq) (*GetModelListRsp, error) {
	geminiReq := &GeminiGetModelListReq{
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	}
	geminiRsp, err := a.client.GetModelList(ctx, geminiReq)
	if err != nil {
		return nil, err
	}

	rsp := &GetModelListRsp{
		Object:        "list",
		NextPageToken: geminiRsp.NextPageToken,
	}
	for _, model := range geminiRsp.Models {
		rsp.Data = append(rsp.Data, struct {
			Id                         string             `json:"id"`
			Name                       string             `json:"name"`
			Object                     string             `json:"object"`
			Created                    int                `json:"created"`
			OwnedBy                    string             `json:"owned_by"`
			Model                      string             `json:"model"`
			ModifiedAt                 time.Time          `json:"modified_at"`
			Size                       int64              `json:"size"`
			Digest                     string             `json:"digest"`
			Details                    OllamaModelDetails `json:"details"`
			BaseModelID                string             `json:"baseModelId"`
			Version                    string             `json:"version"`
			DisplayName                string             `json:"displayName"`
			Description                string             `json:"description"`
			InputTokenLimit            int                `json:"inputTokenLimit"`
			OutputTokenLimit           int                `json:"outputTokenLimit"`
			SupportedGenerationMethods []string           `json:"supportedGenerationMethods"`
			Thinking                   bool               `json:"thinking"`
			Temperature                int                `json:"temperature"`
			MaxTemperature             int                `json:"maxTemperature"`
			TopP                       int                `json:"topP"`
			TopK                       int                `json:"topK"`
		}{
			Id:                         model.Name,
			Name:                       model.Name,
			Object:                     "model",
			OwnedBy:                    "Google",
			BaseModelID:                model.BaseModelID,
			Version:                    model.Version,
			DisplayName:                model.DisplayName,
			Description:                model.Description,
			InputTokenLimit:            model.InputTokenLimit,
			OutputTokenLimit:           model.OutputTokenLimit,
			SupportedGenerationMethods: model.SupportedGenerationMethods,
			Temperature:                anyx.ToInt(model.Temperature),
			TopP:                       anyx.ToInt(model.TopP),
			TopK:                       anyx.ToInt(model.TopK),
		})
	}

	return rsp, nil
}

// SiliconFlowAdapter 实现了 Channeler 接口，用于适配 SiliconFlow 客户端。
type SiliconFlowAdapter struct {
	client *SiliconFlow
}

// NewSiliconFlowAdapter 创建一个新的 SiliconFlow 适配器实例。
func NewSiliconFlowAdapter(channel *aiload.ModelChannel) *SiliconFlowAdapter {
	return &SiliconFlowAdapter{
		client: NewSiliconFlow(channel),
	}
}

// GetModelList 调用 SiliconFlow 客户端的 GetModelList 方法，并将结果转换为统一的 GetModelListRsp 格式。
func (a *SiliconFlowAdapter) GetModelList(ctx context.Context, req *GetModelListReq) (*GetModelListRsp, error) {
	siliconFlowReq := &SiliconFlowGetModelListReq{
		Type:    req.Type,
		SubType: req.SubType,
	}
	siliconFlowRsp, err := a.client.GetModelList(siliconFlowReq)
	if err != nil {
		return nil, err
	}

	rsp := &GetModelListRsp{
		Object: siliconFlowRsp.Object,
	}
	for _, model := range siliconFlowRsp.Data {
		rsp.Data = append(rsp.Data, struct {
			Id                         string             `json:"id"`
			Name                       string             `json:"name"`
			Object                     string             `json:"object"`
			Created                    int                `json:"created"`
			OwnedBy                    string             `json:"owned_by"`
			Model                      string             `json:"model"`
			ModifiedAt                 time.Time          `json:"modified_at"`
			Size                       int64              `json:"size"`
			Digest                     string             `json:"digest"`
			Details                    OllamaModelDetails `json:"details"`
			BaseModelID                string             `json:"baseModelId"`
			Version                    string             `json:"version"`
			DisplayName                string             `json:"displayName"`
			Description                string             `json:"description"`
			InputTokenLimit            int                `json:"inputTokenLimit"`
			OutputTokenLimit           int                `json:"outputTokenLimit"`
			SupportedGenerationMethods []string           `json:"supportedGenerationMethods"`
			Thinking                   bool               `json:"thinking"`
			Temperature                int                `json:"temperature"`
			MaxTemperature             int                `json:"maxTemperature"`
			TopP                       int                `json:"topP"`
			TopK                       int                `json:"topK"`
		}{
			Id:      model.Id,
			Name:    model.Id,
			Object:  model.Object,
			Created: int(model.Created),
			OwnedBy: model.OwnedBy,
		})
	}

	return rsp, nil
}

// OpenAiAdapter 实现了 Channeler 接口，用于适配 OpenAI 客户端。
type OpenAiAdapter struct {
	client *OpenAI
}

// NewOpenAiAdapter 创建一个新的 OpenAI 适配器实例。
func NewOpenAiAdapter(channel *aiload.ModelChannel) *OpenAiAdapter {
	return &OpenAiAdapter{
		client: NewOpenAI(channel),
	}
}

// GetModelList 调用 OpenAI 客户端的 GetModelList 方法，并将结果转换为统一的 GetModelListRsp 格式。
func (a *OpenAiAdapter) GetModelList(ctx context.Context, req *GetModelListReq) (*GetModelListRsp, error) {
	openAiRsp, err := a.client.GetModelList()
	if err != nil {
		return nil, err
	}

	rsp := &GetModelListRsp{
		Object: openAiRsp.Object,
	}
	for _, model := range openAiRsp.Data {
		rsp.Data = append(rsp.Data, struct {
			Id                         string             `json:"id"`
			Name                       string             `json:"name"`
			Object                     string             `json:"object"`
			Created                    int                `json:"created"`
			OwnedBy                    string             `json:"owned_by"`
			Model                      string             `json:"model"`
			ModifiedAt                 time.Time          `json:"modified_at"`
			Size                       int64              `json:"size"`
			Digest                     string             `json:"digest"`
			Details                    OllamaModelDetails `json:"details"`
			BaseModelID                string             `json:"baseModelId"`
			Version                    string             `json:"version"`
			DisplayName                string             `json:"displayName"`
			Description                string             `json:"description"`
			InputTokenLimit            int                `json:"inputTokenLimit"`
			OutputTokenLimit           int                `json:"outputTokenLimit"`
			SupportedGenerationMethods []string           `json:"supportedGenerationMethods"`
			Thinking                   bool               `json:"thinking"`
			Temperature                int                `json:"temperature"`
			MaxTemperature             int                `json:"maxTemperature"`
			TopP                       int                `json:"topP"`
			TopK                       int                `json:"topK"`
		}{
			Id:      model.Id,
			Name:    model.Id,
			Object:  model.Object,
			Created: model.Created,
			OwnedBy: model.OwnedBy,
		})
	}

	return rsp, nil
}

// OpenAiCompatibleAdapter 实现了 Channeler 接口，用于适配与 OpenAI 兼容的客户端。
type OpenAiCompatibleAdapter struct {
	client *OpenAICompatible
}

// NewOpenAiCompatibleAdapter 创建一个新的 OpenAiCompatible 适配器实例。
func NewOpenAiCompatibleAdapter(channel *aiload.ModelChannel) *OpenAiCompatibleAdapter {
	return &OpenAiCompatibleAdapter{
		client: NewOpenAICompatible(channel),
	}
}

// GetModelList 调用 OpenAiCompatible 客户端的 GetModelList 方法，并将结果转换为统一的 GetModelListRsp 格式。
func (a *OpenAiCompatibleAdapter) GetModelList(ctx context.Context, req *GetModelListReq) (*GetModelListRsp, error) {
	openAiCompatRsp, err := a.client.GetModelList()
	if err != nil {
		return nil, err
	}

	rsp := &GetModelListRsp{
		Object: openAiCompatRsp.Object,
	}
	for _, model := range openAiCompatRsp.Data {
		rsp.Data = append(rsp.Data, struct {
			Id                         string             `json:"id"`
			Name                       string             `json:"name"`
			Object                     string             `json:"object"`
			Created                    int                `json:"created"`
			OwnedBy                    string             `json:"owned_by"`
			Model                      string             `json:"model"`
			ModifiedAt                 time.Time          `json:"modified_at"`
			Size                       int64              `json:"size"`
			Digest                     string             `json:"digest"`
			Details                    OllamaModelDetails `json:"details"`
			BaseModelID                string             `json:"baseModelId"`
			Version                    string             `json:"version"`
			DisplayName                string             `json:"displayName"`
			Description                string             `json:"description"`
			InputTokenLimit            int                `json:"inputTokenLimit"`
			OutputTokenLimit           int                `json:"outputTokenLimit"`
			SupportedGenerationMethods []string           `json:"supportedGenerationMethods"`
			Thinking                   bool               `json:"thinking"`
			Temperature                int                `json:"temperature"`
			MaxTemperature             int                `json:"maxTemperature"`
			TopP                       int                `json:"topP"`
			TopK                       int                `json:"topK"`
		}{
			Id:      model.Id,
			Name:    model.Id,
			Object:  model.Object,
			Created: model.Created,
			OwnedBy: model.OwnedBy,
		})
	}

	return rsp, nil
}

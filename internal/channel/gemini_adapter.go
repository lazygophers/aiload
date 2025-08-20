package channel

import (
	"time"

	"github.com/lazygophers/aiload"
	"github.com/lazygophers/utils/anyx"
)

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
func (a *GeminiAdapter) GetModelList() (*GetModelListRsp, error) {
	/*geminiReq := &GeminiGetModelListReq{
		PageSize:  req.PageSize,
		PageToken: req.PageToken,
	}*/
	geminiRsp, err := a.client.GetModelList()
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

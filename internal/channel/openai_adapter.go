package channel

import (
	"context"
	"time"

	"github.com/lazygophers/aiload"
)

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

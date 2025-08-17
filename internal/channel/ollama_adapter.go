package channel

import (
	"context"
	"time"

	"github.com/lazygophers/aiload"
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

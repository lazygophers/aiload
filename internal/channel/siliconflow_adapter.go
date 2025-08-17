package channel

import (
	"context"
	"time"

	"github.com/lazygophers/aiload"
)

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

package channel

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/anyx"
)

type Gemini struct {
	channel *aiload.ModelChannel
}

func NewGemini(channel *aiload.ModelChannel) *Gemini {
	return &Gemini{
		channel: channel,
	}
}

type GeminiGetModelListReq struct {
	PageSize  int
	PageToken string
}

type GeminiModel struct {
	Name                       string   `json:"name"`
	BaseModelId                string   `json:"baseModelId"`
	Version                    string   `json:"version"`
	DisplayName                string   `json:"displayName"`
	Description                string   `json:"description"`
	InputTokenLimit            int      `json:"inputTokenLimit"`
	OutputTokenLimit           int      `json:"outputTokenLimit"`
	SupportedGenerationMethods []string `json:"supportedGenerationMethods"`
	Thinking                   bool     `json:"thinking"`
	Temperature                int      `json:"temperature"`
	MaxTemperature             int      `json:"maxTemperature"`
	TopP                       int      `json:"topP"`
	TopK                       int      `json:"topK"`
}

type GeminiGetModelListRsp struct {
	Models        []GeminiModel `json:"models"`
	NextPageToken string        `json:"nextPageToken"`
}

func (p *Gemini) GetModelList(ctx context.Context, req *GeminiGetModelListReq) (*GeminiGetModelListRsp, error) {
	var rsp GeminiGetModelListRsp
	pageSize := anyx.ToString(req.PageSize)
	if req.PageSize == 0 {
		pageSize = "50"
	}
	if req.PageSize > 1000 {
		pageSize = "1000"
	}

	resp, err := p.GetRequest().SetQueryParams(map[string]string{
		"pageSize":  pageSize,
		"pageToken": req.PageToken,
	}).Get("https://generativelanguage.googleapis.com/v1beta/models")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if resp.IsError() {
		return nil, errors.New(resp.String())
	}

	if err := json.Unmarshal(resp.Body(), &rsp); err != nil {
		return nil, err
	}

	return &rsp, nil
}

func (p *Gemini) GetRequest() *resty.Request {
	return client.R().SetBody(nil)
}

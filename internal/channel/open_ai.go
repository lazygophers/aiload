package channel

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
)

type OpenAI struct {
	channel *aiload.ModelChannel
	baseURL string
}

func NewOpenAI(channel *aiload.ModelChannel) *OpenAI {
	baseURL := "https://api.openai.com"
	if channel.BaseUrl != "" {
		baseURL = channel.BaseUrl
	}
	return &OpenAI{
		channel: channel,
		baseURL: baseURL,
	}
}

type OpenAIGetModelListRsp struct {
	Object string `json:"object"`
	Data   []struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func (p *OpenAI) GetModelList() (*OpenAIGetModelListRsp, error) {
	var rsp OpenAIGetModelListRsp
	resp, err := p.GetRequest().Get(p.baseURL + "/v1/models")
	if err != nil {
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

func (p *OpenAI) GetRequest() *resty.Request {
	// NOTE: The baseURL is prepended to the URL in each request method
	// because resty doesn't allow setting BaseURL on a per-request basis.
	return client.R().
		SetAuthToken(p.channel.Token)
}

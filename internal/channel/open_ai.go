package channel

import (
	"encoding/json"
	"errors"

	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
)

type OpenAI struct {
	channel *aiload.ModelChannel
}

func NewOpenAI(channel *aiload.ModelChannel) *OpenAI {
	return &OpenAI{
		channel: channel,
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
	resp, err := p.GetRequest().Get("https://api.openai.com/v1/models")
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
	return client.R().
		SetAuthToken(p.channel.Token)
}

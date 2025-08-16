package channel

import (
	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
)

type OpenAICompatible struct {
	channel *aiload.ModelChannel
}

func NewOpenAICompatible(channel *aiload.ModelChannel) *OpenAICompatible {
	return &OpenAICompatible{
		channel: channel,
	}
}

type OpenAICompatibleGetModelListRsp struct {
	Data []struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
	Object string `json:"object"`
}

func (p *OpenAICompatible) GetModelList() (*OpenAICompatibleGetModelListRsp, error) {
	var rsp OpenAICompatibleGetModelListRsp
	_, err := p.GetRequest().SetResult(&rsp).Get("models")
	if err != nil {
		return nil, err
	}

	return &rsp, nil
}

func (p *OpenAICompatible) GetRequest() *resty.Request {
	return client.R().
		SetHeader("Authorization", "Bearer "+p.channel.Token)
}

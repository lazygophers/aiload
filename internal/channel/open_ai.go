package channel

import (
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
	_, err := p.GetRequest().SetResult(&rsp).Get("https://api.openai.com/v1/models")
	if err != nil {
		return nil, err
	}
	return &rsp, nil
}

func (p *OpenAI) GetRequest() *resty.Request {
	return client.R().
		SetAuthToken(p.channel.Token)
}

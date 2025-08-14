package channel

import (
	"errors"
	"github.com/go-resty/resty/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
)

var ErrAuthTokenNil = errors.New("auth token is nil")

type SiliconFlow struct {
	channel *aiload.ModelChannel
}

type SiliconFlowGetInfoListRsp struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  bool   `json:"status"`
	Data    struct {
		Id            string `json:"id"`
		Name          string `json:"name"`
		Image         string `json:"image"`
		Email         string `json:"email"`
		IsAdmin       bool   `json:"isAdmin"`
		Balance       string `json:"balance"`
		Status        string `json:"status"`
		Introduction  string `json:"introduction"`
		Role          string `json:"role"`
		ChargeBalance string `json:"chargeBalance"`
		TotalBalance  string `json:"totalBalance"`
	} `json:"data"`
}

type SiliconFlowGetModelListReq struct {
	Type    string
	SubType string
}

type SiliconFlowGetModelListRsp struct {
	Object string `json:"object"`
	Data   []struct {
		Id      string `json:"id"`
		Object  string `json:"object"`
		Created int    `json:"created"`
		OwnedBy string `json:"owned_by"`
	} `json:"data"`
}

func (p *SiliconFlow) GetInfoList() (*SiliconFlowGetInfoListRsp, error) {
	var rsp SiliconFlowGetInfoListRsp
	if p.GetRequestRequired() == false {
		return nil, ErrAuthTokenNil
	}

	_, err := p.GetRequest().SetResult(&rsp).Get("https://api.siliconflow.cn/v1/user/info")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	return &rsp, nil
}

func (p *SiliconFlow) GetModelList(req *SiliconFlowGetModelListReq) (*SiliconFlowGetModelListRsp, error) {
	var rsp SiliconFlowGetModelListRsp
	_, err := p.GetRequest().SetResult(&rsp).SetQueryParams(map[string]string{
		"type":     req.Type,
		"sub_type": req.SubType,
	}).Get("https://api.siliconflow.cn/v1/models")
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	return &rsp, nil
}

func (p *SiliconFlow) GetRequest() *resty.Request {
	return client.R().
		SetAuthToken(p.channel.Token)
}

func (p *SiliconFlow) GetRequestRequired() bool {
	if p.channel.Token == "" {
		return false
	}

	return true
}

func NewSiliconFlow(channel *aiload.ModelChannel) *SiliconFlow {
	p := &SiliconFlow{
		channel: channel,
	}

	return p
}

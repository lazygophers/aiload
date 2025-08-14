package platform

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/log"
)

type SiliconFlow struct {
	obj *aiload.ModelChannel
}

type SiliconFlowGetModelListRsp struct {
}

func (s *SiliconFlow) GetModelList() (*SiliconFlowGetModelListRsp, error) {
	var rsp SiliconFlowGetModelListRsp
	_, err := client.R().
		SetHeaders(map[string]string{
			"Authorization": "Bearer " + obj.Token,
		}).
		SetResult(&rsp).
		Get("https://api.siliconflow.cn/v1/models")
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func NewSiliconFlow(obj *aiload.ModelChannel) *SiliconFlow {
	return &SiliconFlow{
		obj: obj,
	}
}

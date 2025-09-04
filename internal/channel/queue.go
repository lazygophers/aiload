package channel

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/memqueue"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/utils/anyx"
	"github.com/lazygophers/utils/candy"
)

func LoadQueue() {
	state.MqUpdateModel.Process(func(msg *memqueue.Message[*aiload.MqTaskId]) (*memqueue.Consume, error) {
		var rsp memqueue.Consume

		obj, err := state.Channel.NewScoop().
			Where(aiload.DbId, msg.Data.Id).
			In(aiload.DbState, []int32{
				int32(aiload.State_Enable),
				int32(aiload.State_AutoDisable),
			}).
			Where(aiload.DbEnableAutoUpdateModel, true).
			First()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = UpdateModelList(obj)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return &rsp, nil
	})

	state.MqCheckChannelModel.Process(func(msg *memqueue.Message[*aiload.MqTaskId]) (*memqueue.Consume, error) {
		var rsp memqueue.Consume

		channelModel, err := state.ChannelModel.NewScoop().
			Where(aiload.DbId, msg.Data.Id).
			In(aiload.DbState, []int32{
				int32(aiload.State_Enable),
				int32(aiload.State_AutoDisable),
			}).
			First()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		channel, err := state.Channel.NewScoop().
			Where(aiload.DbId, channelModel.ChannelId).
			In(aiload.DbState, []int32{
				int32(aiload.State_Enable),
				int32(aiload.State_AutoDisable),
			}).
			Where(aiload.DbEnableAutoCheckModel, true).
			First()
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		err = CheckChannelModel(channel, channelModel)
		if err != nil {
			log.Errorf("err:%v", err)
			return nil, err
		}

		return &rsp, nil
	})
}

func UpdateModelList(obj *aiload.ModelChannel) error {
	channel, err := NewChannel(obj)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	resp, err := channel.GetModelList()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = state.CommitOrRollback(func(tx *db.Scoop) error {
		var oldModelList []string
		newModelList := anyx.PluckString(resp.Data, "Id")

		{
			modelList, err := state.ChannelModel.NewScoop(tx).
				Select(aiload.DbModel).
				Where(aiload.DbChannelId, obj.Id).
				Find()
			if err != nil {
				log.Errorf("err:%s", err)
				return err
			}

			oldModelList = anyx.PluckString(modelList, "Model")
		}

		added, removed := candy.Diff(oldModelList, newModelList)

		if len(removed) > 0 {
			err = state.ChannelModel.NewScoop(tx).
				Where(aiload.DbChannelId, obj.Id).
				In(aiload.DbModel, removed).
				Delete().
				Error
			if err != nil {
				log.Errorf("err:%s", err)
				return err
			}
		}

		if len(added) > 0 {
			err = state.ChannelModel.NewScoop(tx).
				CreateInBatches(candy.Map(added, func(modelName string) *aiload.ModelChannelModel {
					return &aiload.ModelChannelModel{
						ChannelId: obj.Id,
						Model:     modelName,
					}
				}), 100).
				Error
			if err != nil {
				log.Errorf("err:%s", err)
				return err
			}
		}

		return nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

func CheckChannelModel(obj *aiload.ModelChannel, modelObj *aiload.ModelChannelModel) error {
	channel, err := NewChannel(obj)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	// todo: 首字延时、关键词匹配过滤、存活判断、embedc处理、thinking 处理 （https://github.com/QuantumNous/new-api/blob/main/controller/channel-test.go）

	err = channel.ChatCompletionsAsync(&ChatCompletionReq{
		Model:  modelObj.Model,
		Stream: true,
		Messages: []*ChatMessage{
			{
				Role:    "user",
				Content: "hi",
			},
		},
	}, func(item *ChatCompletionRspItem) error {
		return nil
	})
	if err != nil {
		log.Errorf("err:%v", err)
	} else {

	}

	return err
}

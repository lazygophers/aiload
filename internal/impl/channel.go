package impl

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/channel"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/lrpc/middleware/xerror"
	"github.com/lazygophers/utils/anyx"
	"github.com/lazygophers/utils/candy"
)

func SetChannelAdmin(ctx *lrpc.Ctx, req *aiload.SetChannelAdminReq) (*aiload.SetChannelAdminRsp, error) {
	var rsp aiload.SetChannelAdminRsp

	//goland:noinspection GoVetCopyLock
	channel := req.Channel

	err := state.CommitOrRollback(func(tx *db.Scoop) error {
		err := state.Channel.
			NewScoop(tx).
			Where("id = ?", channel.Id).
			Updates(channel).
			Error
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		var oldModelList []string
		newModelList := channel.ModelList

		{
			modelList, err := state.ChannelAccess.NewScoop(tx).
				Select(aiload.DbModel).
				Where(aiload.DbChannelId, channel.Id).
				Find()
			if err != nil {
				log.Errorf("err:%s", err)
				return err
			}

			oldModelList = anyx.PluckString(modelList, "Model")
		}

		added, removed := candy.Diff(oldModelList, newModelList)

		if len(removed) > 0 {
			err = state.ChannelAccess.NewScoop(tx).
				Where(aiload.DbChannelId, channel.Id).
				In(aiload.DbModel, removed).
				Delete().
				Error
			if err != nil {
				log.Errorf("err:%s", err)
				return err
			}
		}

		if len(added) > 0 {
			err = state.ChannelAccess.NewScoop(tx).
				CreateInBatches(candy.Map(added, func(modelName string) *aiload.ModelChannelAccess {
					return &aiload.ModelChannelAccess{
						ChannelId: channel.Id,
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
		log.Errorf("err:%s", err)
		return nil, err
	}

	rsp.Channel, err = state.Channel.
		NewScoop().
		Where("id = ?", channel.Id).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func DelChannelAdmin(ctx *lrpc.Ctx, req *aiload.DelChannelAdminReq) (*aiload.DelChannelAdminRsp, error) {
	var rsp aiload.DelChannelAdminRsp

	err := state.Channel.
		NewScoop().
		Where("id = ?", req.Id).
		Delete().
		Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func AddChannelAdmin(ctx *lrpc.Ctx, req *aiload.AddChannelAdminReq) (*aiload.AddChannelAdminRsp, error) {
	var rsp aiload.AddChannelAdminRsp

	//goland:noinspection GoVetCopyLock
	channel := *req.Channel
	channel.Id = 0

	err := state.CommitOrRollback(func(tx *db.Scoop) error {
		err := state.Channel.
			NewScoop(tx).
			Create(&channel)
		if err != nil {
			log.Errorf("err:%v", err)
			return err
		}

		err = state.ChannelAccess.NewScoop(tx).
			CreateInBatches(candy.Map(channel.ModelList, func(modelName string) *aiload.ModelChannelAccess {
				return &aiload.ModelChannelAccess{
					ChannelId: channel.Id,
					Model:     modelName,
				}
			}), 100).
			Error
		if err != nil {
			log.Errorf("err:%s", err)
			return err
		}

		return nil
	})
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	rsp.Channel = &channel

	return &rsp, nil
}

func GetChannelAdmin(ctx *lrpc.Ctx, req *aiload.GetChannelAdminReq) (*aiload.GetChannelAdminRsp, error) {
	var rsp aiload.GetChannelAdminRsp

	channel, err := state.Channel.
		NewScoop().
		Where("id = ?", req.Id).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Channel = channel

	return &rsp, nil
}

func ListChannelAdmin(ctx *lrpc.Ctx, req *aiload.ListChannelAdminReq) (*aiload.ListChannelAdminRsp, error) {
	var rsp aiload.ListChannelAdminRsp

	scoop := state.Channel.NewScoop()

	err := req.ListOption.Processor().
		Uint32(int32(aiload.ListChannelAdminReq_ListOptionPlatform), func(value uint32) error {
			scoop.Where(aiload.DbPlatform, value)
			return nil
		}).
		String(int32(aiload.ListChannelAdminReq_ListOptionName), func(value string) error {
			scoop.Like(aiload.DbName, value)
			return nil
		}).
		String(int32(aiload.ListChannelAdminReq_ListOptionToken), func(value string) error {
			scoop.Like(aiload.DbToken, value)
			return nil
		}).
		String(int32(aiload.ListChannelAdminReq_ListOptionBaseUrl), func(value string) error {
			scoop.Like(aiload.DbBaseUrl, value)
			return nil
		}).
		Process()
	if err != nil {
		if xerror.CheckCode(err, xerror.ErrNoData) {
			rsp.Paginate = req.ListOption.Paginate()
			return &rsp, nil
		}

		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.Paginate, rsp.List, err = scoop.FindByPage(req.ListOption)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}
func GetModelListAdmin(ctx *lrpc.Ctx, req *aiload.GetModelListAdminReq) (*aiload.GetModelListAdminRsp, error) {
	var rsp aiload.GetModelListAdminRsp

	channeler, err := channel.NewChannel(req.Channel)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	var getModelList *channel.GetModelListRsp
	getModelList, err = channeler.GetModelList()
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	rsp.ModelList = anyx.PluckString(getModelList.Data, "Model")
	/*candy.Map(getModelList.Data, func(model ) string {
		return &model.Model
	})*/

	return &rsp, nil
}

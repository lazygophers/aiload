package impl

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/xerror"
)

func SetChannelAdmin(ctx *lrpc.Ctx, req *aiload.SetChannelAdminReq) (*aiload.SetChannelAdminRsp, error) {
	var rsp aiload.SetChannelAdminRsp

	//goland:noinspection GoVetCopyLock
	channel := req.Channel

	err := state.Channel.
		NewScoop().
		Where("id = ?", channel.Id).
		Updates(channel).
		Error
	if err != nil {
		log.Errorf("err:%v", err)
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

	err := state.Channel.
		NewScoop().
		Create(&channel)
	if err != nil {
		log.Errorf("err:%v", err)
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

	return &rsp, nil
}

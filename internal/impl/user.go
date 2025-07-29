package impl

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/xerror"
)

func GetUser(ctx *fiber.Ctx, req *aiload.GetUserReq) (*aiload.GetUserRsp, error) {
	var rsp aiload.GetUserRsp

	user, err := state.User.
		NewScoop().
		Where("id = ?", req.Id).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.User = user

	return &rsp, nil
}

func GetUserAdmin(ctx *fiber.Ctx, req *aiload.GetUserAdminReq) (*aiload.GetUserAdminRsp, error) {
	var rsp aiload.GetUserAdminRsp

	user, err := state.User.
		NewScoop().
		Where("id = ?", req.Id).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.User = user

	return &rsp, nil
}

func ListUserAdmin(ctx *fiber.Ctx, req *aiload.ListUserAdminReq) (*aiload.ListUserAdminRsp, error) {
	var rsp aiload.ListUserAdminRsp

	scoop := state.User.NewScoop()

	err := req.ListOption.Processor().
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

func SetUser(ctx *fiber.Ctx, req *aiload.SetUserReq) (*aiload.SetUserRsp, error) {
	var rsp aiload.SetUserRsp

	//goland:noinspection GoVetCopyLock
	user := req.User

	err := state.User.
		NewScoop().
		Where("id = ?", user.Id).
		Updates(user).
		Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.User, err = state.User.
		NewScoop().
		Where("id = ?", user.Id).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func SetUserAdmin(ctx *fiber.Ctx, req *aiload.SetUserAdminReq) (*aiload.SetUserAdminRsp, error) {
	var rsp aiload.SetUserAdminRsp

	//goland:noinspection GoVetCopyLock
	user := req.User

	err := state.User.
		NewScoop().
		Where("id = ?", user.Id).
		Updates(user).
		Error
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.User, err = state.User.
		NewScoop().
		Where("id = ?", user.Id).
		First()
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &rsp, nil
}

func DelUserAdmin(ctx *fiber.Ctx, req *aiload.DelUserAdminReq) (*aiload.DelUserAdminRsp, error) {
	var rsp aiload.DelUserAdminRsp

	err := state.User.
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

func AddUserAdmin(ctx *fiber.Ctx, req *aiload.AddUserAdminReq) (*aiload.AddUserAdminRsp, error) {
	var rsp aiload.AddUserAdminRsp

	//goland:noinspection GoVetCopyLock
	user := *req.User
	user.Id = 0

	err := state.User.
		NewScoop().
		Create(&user)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}

	rsp.User = &user

	return &rsp, nil
}

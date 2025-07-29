package impl

import (
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc"
	"github.com/lazygophers/lrpc/middleware/xerror"
)

func Login(ctx *fiber.Ctx, req *aiload.LoginReq) (*aiload.LoginRsp, error) {
	var rsp aiload.LoginRsp

	return &rsp, nil
}

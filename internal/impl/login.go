package impl

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/lazygophers/aiload"
	"github.com/lazygophers/aiload/internal/state"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/cryptox"
	"github.com/lazygophers/utils/xtime"
	"time"
)

func Login(ctx *fiber.Ctx, req *aiload.LoginReq) (*aiload.LoginRsp, error) {
	var rsp aiload.LoginRsp
	user, err := state.User.
		NewScoop().
		Where(aiload.DbUsername, req.Username).
		First()
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	if user.Password != req.Password {
		log.Errorf("err:%s", "password error")
		return nil, err
	}

	rsp.Token = "us-" + cryptox.UUID()

	err = state.Cache().SetEx(fmt.Sprintf(state.CacheKeySession, rsp.Token), user, xtime.Day)
	if err != nil {
		log.Errorf("err:%s", err)
		return nil, err
	}

	rsp.User = user
	rsp.ExpiresAt = time.Now().Add(xtime.Day).Unix()

	return &rsp, nil
}

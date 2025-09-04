package core

import (
	"github.com/lazygophers/aiload/internal/api"
	"github.com/lazygophers/aiload/internal/channel"
	"github.com/lazygophers/log"
	"github.com/lazygophers/utils/runtime"
)

func Load() (err error) {
	channel.Load()

	err = api.Listen()
	if err != nil {
		log.Errorf("err:%s", err)
		return err
	}

	runtime.WaitExit()

	return
}

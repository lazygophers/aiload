package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/i18n"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/utils/app"
	"github.com/pterm/pterm"
)

type state struct {
	Config *Config

	// NOTE: Please fill in the state below
	I18n *i18n.I18n
}

var State = new(state)

func Load() (err error) {
	log.SetPrefixMsg(app.Name)

	err = LoadI18n()
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	err = LoadConfig()
	if err != nil {
		log.Errorf("err:%v", err)
		pterm.Debug.Println(err)
		return err
	}

	db.DefaultDriver = State.Config.Db.Type

	err = ConnectDatabase()
	if err != nil {
		log.Errorf("err:%v", err)
		pterm.Debug.Println(err)
		return err
	}

	// 清理一下旧的表结构
	err = ConnectCache()
	if err != nil {
		log.Errorf("err:%v", err)
		pterm.Debug.Println(err)
		return err
	}

	return nil
}

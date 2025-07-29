package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/cache"
	"github.com/lazygophers/lrpc/middleware/storage/db"
	"github.com/lazygophers/utils"
	"github.com/lazygophers/utils/config"
)

type Config struct {
	Db    *db.Config    `json:"db,omitempty" yaml:"db,omitempty" toml:"db,omitempty" validate:"required"`
	Cache *cache.Config `json:"cache,omitempty" yaml:"cache,omitempty" toml:"cache,omitempty" validate:"required"`
}

func (p *Config) apply() {
}

func LoadConfig() (err error) {
	State.Config = new(Config)

	_ = config.LoadConfigSkipValidate(State.Config)

	State.Config.apply()

	err = utils.Validate(State.Config)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	return nil
}

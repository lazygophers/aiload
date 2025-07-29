package state

import (
	"github.com/lazygophers/log"
	"github.com/lazygophers/lrpc/middleware/storage/cache"
	"github.com/lazygophers/utils/app"
)

var (
	_cache cache.Cache
)

const (
	CacheKeySession = "session:%s"
)

func ConnectCache() (err error) {
	log.Info("try init cache")
	_cache, err = cache.New(State.Config.Cache)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	_cache.SetPrefix(app.Name + ":")

	return nil
}

func Cache() cache.Cache {
	return _cache
}

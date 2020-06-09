package cache

import (
	log "adeia-api/internal/utils/logger"
	"errors"
	"sync"

	"github.com/mediocregopher/radix/v3"
	config "github.com/spf13/viper"
)

type Cache struct {
	*radix.Pool
}

var (
	pool *Cache
	initCache = new(sync.Once)
)

func Init() error {
	err := errors.New("")

	initCache.Do(func() {
		err = nil
		p, e := radix.NewPool(
			config.GetString("cache.network"),
			config.GetString("cache.addr"),
			config.GetInt("cache.connsize"),
		)
		if e != nil {
			err = e
			return
		}
		pool = &Cache{p}
	})

	return err
}

func Get(rcv interface{}, args ...string) error {
	return do(radix.Cmd(rcv, "GET", args...))
}

func Set(key string, value string) error {
	return do(radix.Cmd(nil, "SET", key, value))
}

func do(cmd radix.CmdAction) error {
	if pool == nil {
		msg := "cache pool not initialized"
		log.Warn(msg)
		return errors.New(msg)
	}
	return pool.Do(cmd)
}

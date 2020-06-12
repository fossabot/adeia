package cache

import (
	log "adeia-api/internal/utils/logger"
	"errors"
	"strconv"
	"sync"

	"github.com/mediocregopher/radix/v3"
	config "github.com/spf13/viper"
)

// Cache represents the cache connection instance.
type Cache struct {
	*radix.Pool
}

var (
	pool      *Cache
	initCache = new(sync.Once)
)

// Init creates a new cache connection.
func Init() error {
	err := errors.New("")

	initCache.Do(func() {
		err = nil
		addr := config.GetString("cache.host") + ":" + config.GetString("cache.port")
		p, e := radix.NewPool(
			config.GetString("cache.network"),
			addr,
			config.GetInt("cache.connsize"),
		)
		if e != nil {
			err = e
			return
		}

		// ping
		pong := ""
		e = p.Do(radix.Cmd(&pong, "PING"))
		if e != nil {
			err = e
			return
		} else if pong != "PONG" {
			err = errors.New("Expected: PONG. Received: " + pong)
			return
		}

		pool = &Cache{p}
	})

	return err
}

// Close closes the cache connection.
func Close() error {
	return pool.Close()
}

// Get gets the value of the specified key.
func Get(rcv interface{}, key string) error {
	return do(radix.Cmd(rcv, "GET", key))
}

// SetWithExpiry sets the provided key:value pair with specified seconds of TTL.
func SetWithExpiry(key string, value string, seconds int) error {
	return do(radix.Cmd(nil, "SET", key, value, "EX", strconv.Itoa(seconds)))
}

// Set sets the provided key:value pair.
func Set(key string, value string) error {
	return do(radix.Cmd(nil, "SET", key, value))
}

// Delete deletes the list of keys.
func Delete(keys ...string) error {
	return do(radix.Cmd(nil, "DEL", keys...))
}

// do is a wrapper on the pool.Do func. All cache functions hit this method.
// This enables a fall-through to the database, when cache is unavailable.
func do(cmd radix.CmdAction) error {
	if pool == nil {
		msg := "cache pool not initialized"
		log.Warn(msg)
		return errors.New(msg)
	}
	return pool.Do(cmd)
}

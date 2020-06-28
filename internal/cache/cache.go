package cache

import (
	"errors"
	"fmt"
	"strconv"

	"adeia-api/internal/util/constants"
	"adeia-api/internal/util/log"

	"github.com/mediocregopher/radix/v3"
	config "github.com/spf13/viper"
)

// Cache represents the funcs required to access a key-value store (like Redis).
type Cache interface {
	Get(rcv interface{}, key string) error
	Set(key string, value string) error
	SetWithExpiry(key, value string, seconds uint64) error
	StoreSession(sessID string, userID int) error
	GetSession(rcv interface{}, sessID string) error
	Delete(keys ...string) error
	Close() error
}

// RedisCache represents the cache connection instance.
type RedisCache struct {
	*radix.Pool
}

// New creates a new cache connection instance.
func New() (Cache, error) {
	// TODO: add cache auth
	p, err := radix.NewPool(
		config.GetString("cache.network"),
		config.GetString("cache.host")+":"+config.GetString("cache.port"),
		config.GetInt("cache.connsize"),
	)
	if err != nil {
		return &RedisCache{}, err
	}

	// check connection
	if err := pingCheck(p); err != nil {
		return nil, fmt.Errorf("cannot ping redis instance: %v", err)
	}
	return &RedisCache{p}, nil
}

func pingCheck(p *radix.Pool) error {
	return p.Do(radix.Cmd(nil, "PING"))
}

// Get gets the value of the specified key.
func (r *RedisCache) Get(rcv interface{}, key string) error {
	return r.do(radix.Cmd(rcv, "GET", key))
}

// Set sets the provided key:value pair.
func (r *RedisCache) Set(key string, value string) error {
	return r.do(radix.Cmd(nil, "SET", key, value))
}

// Delete deletes the list of keys.
func (r *RedisCache) Delete(keys ...string) error {
	return r.do(radix.Cmd(nil, "DEL", keys...))
}

// SetWithExpiry sets the provided key:value pair with specified seconds of TTL.
func (r *RedisCache) SetWithExpiry(key, value string, seconds uint64) error {
	return r.do(radix.Cmd(nil, "SET", key, value, "EX", strconv.FormatUint(seconds, 10)))
}

// Close closes the connection pool.
func (r *RedisCache) Close() error {
	if !r.isAvailable() {
		return nil
	}
	return r.Pool.Close()
}

// StoreSession stores a sessID:userID pair.
func (r *RedisCache) StoreSession(sessID string, userID int) error {
	sessKey := buildSessionKey(sessID)
	return r.SetWithExpiry(sessKey, strconv.Itoa(userID), constants.SessionExpiry)
}

// GetSession returns the userID for the provided sessID.
func (r *RedisCache) GetSession(userID interface{}, sessID string) error {
	sessKey := buildSessionKey(sessID)
	return r.Get(userID, sessKey)
}

func buildSessionKey(sessID string) string {
	return "session:" + sessID
}

func (r *RedisCache) isAvailable() bool {
	return r.Pool != nil
}

// do is a wrapper on the pool.Do func. All cache functions hit this method.
// This enables a fall-through to the database, when cache is unavailable.
func (r *RedisCache) do(cmd radix.CmdAction) error {
	if !r.isAvailable() {
		msg := "cache unavailable"
		log.Debug(msg)
		return errors.New(msg)
	}
	return r.Pool.Do(cmd)
}

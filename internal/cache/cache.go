package cache

import (
	"fmt"
	"strconv"

	"github.com/mediocregopher/radix/v3"
	config "github.com/spf13/viper"
)

// Cache represents the funcs required to access a key-value store (like Redis).
type Cache interface {
	Get(rcv interface{}, key string) error
	Set(key string, value string) error
	Expire(key string, seconds int) error
	SetWithExpiry(key, value string, seconds int) error
	Delete(keys ...string) error
	Close() error
	GetInstance() *RedisCache
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

// GetInstance returns the underlying instance.
// TODO: handle this differently because the interface is tied to its implementation.
func (r *RedisCache) GetInstance() *RedisCache {
	return r
}

// Get gets the value of the specified key.
func (r *RedisCache) Get(rcv interface{}, key string) error {
	return r.Do(radix.Cmd(rcv, "GET", key))
}

// Set sets the provided key:value pair.
func (r *RedisCache) Set(key string, value string) error {
	return r.Do(radix.Cmd(nil, "SET", key, value))
}

// Delete deletes the list of keys.
func (r *RedisCache) Delete(keys ...string) error {
	return r.Do(radix.Cmd(nil, "DEL", keys...))
}

// SetWithExpiry sets the provided key:value pair with specified seconds of TTL.
func (r *RedisCache) SetWithExpiry(key, value string, seconds int) error {
	return r.Do(radix.Cmd(nil, "SET", key, value, "EX", strconv.Itoa(seconds)))
}

// Expire sets the expiry for a given key.
func (r *RedisCache) Expire(key string, seconds int) error {
	return r.Do(radix.Cmd(nil, "EXPIRE", key, strconv.Itoa(seconds)))
}

// Close closes the connection pool.
func (r *RedisCache) Close() error {
	return r.Pool.Close()
}

func buildKey(resource, id string, fields ...string) string {
	key := resource + ":" + id
	for _, field := range fields {
		key += ":" + field
	}
	return key
}

package cache

import (
	"strconv"

	"adeia/internal/config"

	"github.com/mediocregopher/radix/v3"
)

// Redis represents the cache connection instance.
type Redis struct {
	*radix.Pool
}

// New creates a new cache connection instance.
func New(conf *config.CacheConfig) (*Redis, error) {
	// TODO: add cache auth
	p, err := radix.NewPool(
		conf.Network,
		conf.Host+":"+strconv.Itoa(conf.Port),
		conf.ConnSize,
	)
	if err != nil {
		return &Redis{}, err
	}

	return &Redis{p}, nil
}

// Get gets the value of the specified key.
func (r *Redis) Get(rcv interface{}, key string) error {
	return r.Do(radix.Cmd(rcv, "GET", key))
}

// Set sets the provided key:value pair.
func (r *Redis) Set(key string, value string) error {
	return r.Do(radix.Cmd(nil, "SET", key, value))
}

// Delete deletes the list of keys.
func (r *Redis) Delete(keys ...string) error {
	return r.Do(radix.Cmd(nil, "DEL", keys...))
}

// SetWithExpiry sets the provided key:value pair with specified seconds of TTL.
func (r *Redis) SetWithExpiry(key, value string, seconds int) error {
	return r.Do(radix.Cmd(nil, "SET", key, value, "EX", strconv.Itoa(seconds)))
}

// Expire sets the expiry for a given key.
func (r *Redis) Expire(key string, seconds int) error {
	return r.Do(radix.Cmd(nil, "EXPIRE", key, strconv.Itoa(seconds)))
}

// Close closes the connection pool.
func (r *Redis) Close() error {
	return r.Pool.Close()
}

func buildKey(resource, id string, fields ...string) string {
	key := resource + ":" + id
	for _, field := range fields {
		key += ":" + field
	}
	return key
}

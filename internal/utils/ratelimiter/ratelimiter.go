package ratelimiter

import (
	config "github.com/spf13/viper"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// A visitor represents a visitor that makes a request. Each visitor has their own
// rate-limiter and lastSeen. The visitor is identified using the request IP.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter represents a rate-limiter.
type RateLimiter struct {
	// map of visitors
	visitors map[string]*visitor

	// mutex for concurrent access
	mu *sync.RWMutex

	// no. of requests allowed per second per visitor
	r rate.Limit

	// no. of burst requests allowed at a time
	b int
}

// NewRateLimiter returns a new RateLimiter with the specified params.
func NewRateLimiter(r rate.Limit, b int) *RateLimiter {
	rl := &RateLimiter{
		visitors: make(map[string]*visitor),
		mu:       &sync.RWMutex{},
		r:        r,
		b:        b,
	}

	// start the cleanup routine
	rl.RunCleanup()
	return rl
}

// AddIP adds the provided ip to the visitors map, marking lastSeen as time.Now().
func (rl *RateLimiter) AddIP(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	limiter := rate.NewLimiter(rl.r, rl.b)
	rl.visitors[ip] = &visitor{limiter, time.Now()}
	return limiter
}

// GetLimiter returns the *rate.Limiter for the provided ip.
func (rl *RateLimiter) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()

	v, exists := rl.visitors[ip]
	if !exists {
		rl.mu.Unlock()
		// ip does not exist, so add it to the map
		return rl.AddIP(ip)
	}

	// ip exists, update lastSeen
	v.lastSeen = time.Now()
	rl.mu.Unlock()
	return v.limiter
}

// RunCleanup is goroutine that regularly removes visitors who've not made a
// request for a period of time, from the map.
func (rl *RateLimiter) RunCleanup() {
	go func() {
		for {
			window := time.Duration(config.GetInt("server.ratelimit_window")) * time.Second
			time.Sleep(window)

			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > window {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}

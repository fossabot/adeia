package ratelimiter

import (
	"sync"
	"time"

	"golang.org/x/time/rate"
)

// visitor represents a visitor that makes a request. Each visitor has their own
// rate-limiter and lastSeen. The visitor is identified using the request IP.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

// RateLimiter is the interface for IP-based rate-limiting.
type RateLimiter interface {
	GetLimiter(ip string) *rate.Limiter
	runCleanup()
}

// Impl represents a RateLimiter implementation.
type Impl struct {
	// map of visitors
	visitors map[string]*visitor

	// mutex for concurrent access
	mu *sync.RWMutex

	// no. of requests allowed per second per visitor
	refillRate rate.Limit

	// no. of burst requests allowed at a time
	bucketSize int

	// how long a visitor can exist in the visitors map without making a request
	cleanupWindow time.Duration
}

// New returns a new Impl with the specified params.
func New(refillRate float64, bucketSize int, duration time.Duration) *Impl {
	rl := &Impl{
		visitors:      make(map[string]*visitor),
		mu:            &sync.RWMutex{},
		refillRate:    rate.Limit(refillRate),
		bucketSize:    bucketSize,
		cleanupWindow: duration,
	}

	// start the cleanup routine
	rl.runCleanup()
	return rl
}

// GetLimiter returns the *rate.Limiter for the provided ip.
func (rl *Impl) GetLimiter(ip string) *rate.Limiter {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	v, exists := rl.visitors[ip]
	if exists {
		// ip exists, update lastSeen
		v.lastSeen = time.Now()
		return v.limiter
	}

	// ip does not exist, so add it to visitors
	rl.visitors[ip] = &visitor{
		limiter:  rate.NewLimiter(rl.refillRate, rl.bucketSize),
		lastSeen: time.Now(),
	}
	return rl.visitors[ip].limiter
}

// runCleanup is goroutine that regularly removes visitors who've not made a
// request for a `window` amount of time, from the map.
func (rl *Impl) runCleanup() {
	// TODO: optimize this
	go func() {
		for {
			time.Sleep(rl.cleanupWindow)

			rl.mu.Lock()
			for ip, v := range rl.visitors {
				if time.Since(v.lastSeen) > rl.cleanupWindow {
					delete(rl.visitors, ip)
				}
			}
			rl.mu.Unlock()
		}
	}()
}

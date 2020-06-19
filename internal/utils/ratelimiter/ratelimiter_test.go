package ratelimiter

import (
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/time/rate"
)

func TestNewRateLimiter(t *testing.T) {
	tests := []struct {
		name string
		r    rate.Limit
		b    int
		d    time.Duration
		want *Impl
		msg  string
	}{
		{
			name: "normal",
			r:    1,
			b:    5,
			d:    time.Duration(5) * time.Second,
			want: &Impl{
				r:             1,
				b:             5,
				mu:            &sync.RWMutex{},
				cleanupWindow: 5 * time.Second,
				visitors:      make(map[string]*visitor),
			},
			msg: "should return a rate limiter with passed-in rate and burst",
		},
		{
			name: "inf rate",
			r:    rate.Inf,
			b:    5,
			d:    5 * time.Second,
			want: &Impl{
				r:             rate.Inf,
				b:             5,
				mu:            &sync.RWMutex{},
				cleanupWindow: 5 * time.Second,
				visitors:      make(map[string]*visitor),
			},
			msg: "should return a rate limiter with passed-in rate and burst",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := New(tc.r, tc.b, tc.d)
			assert.Equal(t, got, tc.want, tc.msg)
		})
	}
}

func toMillisecond(t time.Time) int64 {
	return t.UnixNano() / int64(time.Millisecond)
}

func TestImpl_GetLimiter(t *testing.T) {
	ip := "127.0.0.1"
	rl := New(1, 5, 5*time.Second)

	start := time.Now()
	_ = rl.GetLimiter(ip)
	sec := time.Since(start)

	rl.mu.Lock()
	assert.Equal(t, len(rl.visitors), 1, "should add ip to visitor map")
	assert.Equal(
		t,
		toMillisecond(rl.visitors[ip].lastSeen),
		toMillisecond(start.Add(-1*sec)),
		"should store proper lastSeen time",
	)
	assert.Contains(t, rl.visitors, ip, "ip should be present in the visitor map")
	rl.mu.Unlock()

	start = time.Now()
	_ = rl.GetLimiter(ip)
	sec = time.Since(start)

	rl.mu.Lock()
	assert.Equal(t, len(rl.visitors), 1, "should not new entry to visitor map if ip exists")
	assert.Equal(
		t,
		toMillisecond(rl.visitors[ip].lastSeen),
		toMillisecond(start.Add(-1*sec)),
		"should store proper lastSeen time",
	)
	rl.mu.Unlock()

	_ = rl.GetLimiter("1.1.1.2")
	rl.mu.Lock()
	assert.Equal(t, len(rl.visitors), 2, "should add ip to visitor map")
	rl.mu.Unlock()
}

func TestImpl_RunCleanup(t *testing.T) {
	duration := 5
	ip := "127.0.0.1"

	rl := New(1, 5, time.Duration(duration)*time.Second)
	_ = rl.GetLimiter(ip)

	time.Sleep(time.Duration(duration+1) * time.Second)
	rl.mu.Lock()
	assert.Equal(t, 0, len(rl.visitors), "ip should not exist after window")
	assert.NotContains(t, rl.visitors, ip, "ip should not exist after window")
	rl.mu.Unlock()
}

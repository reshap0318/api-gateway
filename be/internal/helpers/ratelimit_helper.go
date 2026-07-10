package helpers

import (
	"sync"
	"time"
)

// nowFunc allows mocking time.Now() for testing
var nowFunc = time.Now

// Now returns current time (wrapper for testing)
func Now() time.Time {
	return nowFunc()
}

// RateLimitInfo contains rate limit information for response headers
type RateLimitInfo struct {
	Limit     int   // Maximum requests allowed per window
	Remaining int   // Remaining requests in current window
	Reset     int64 // Unix timestamp when the limit resets
	Allowed   bool  // Whether the request is allowed
}

// bucket represents a token bucket for a single client (thread-safe)
type bucket struct {
	mu         sync.Mutex
	tokens     int
	lastRefill time.Time
	windowSecs int // window this bucket was last sized for, used by the reaper to detect staleness
}

// RateLimiter implements a token bucket rate limiter
type RateLimiter struct {
	limit      int
	windowSecs int
	buckets    sync.Map // map[string]*bucket
}

// bucketReapInterval controls how often stale buckets are swept from memory. Buckets whose
// window has been expired for a full window's length (i.e. no request has refreshed them) are
// evicted, otherwise every distinct rate-limit key (e.g. per route+IP) accumulates forever.
const bucketReapInterval = 5 * time.Minute

// NewRateLimiter creates a new rate limiter and starts its background bucket reaper. The
// reaper runs for the lifetime of the process since RateLimiter is a process-wide singleton.
func NewRateLimiter(limit, windowSecs int) *RateLimiter {
	rl := &RateLimiter{
		limit:      limit,
		windowSecs: windowSecs,
	}
	go rl.reapLoop()
	return rl
}

// reapLoop periodically evicts buckets that have been idle for at least one full window past
// their last refill, so a client hitting a route once doesn't leave a bucket in memory forever.
func (rl *RateLimiter) reapLoop() {
	ticker := time.NewTicker(bucketReapInterval)
	defer ticker.Stop()
	for range ticker.C {
		now := time.Now()
		rl.buckets.Range(func(key, value interface{}) bool {
			b := value.(*bucket)
			b.mu.Lock()
			stale := now.Sub(b.lastRefill) >= 2*time.Duration(b.windowSecs)*time.Second
			b.mu.Unlock()
			if stale {
				rl.buckets.Delete(key)
			}
			return true
		})
	}
}

// Allow checks if a request from the given IP is allowed (using the limiter's own
// fixed limit/window) and returns rate limit info.
func (rl *RateLimiter) Allow(ip string) *RateLimitInfo {
	return rl.AllowWithLimit(ip, rl.limit, rl.windowSecs)
}

// AllowWithLimit checks if a request identified by an arbitrary key is allowed against a
// caller-supplied limit/window, independent of the limiter's own default. Used by the
// Dynamic Proxy Engine to enforce per-Route/per-Service rate limits (FSD §2.15) — each
// distinct key gets its own token bucket, so overriding the limit for one route/service
// does not affect any other key's bucket.
func (rl *RateLimiter) AllowWithLimit(key string, limit, windowSecs int) *RateLimitInfo {
	now := time.Now()
	resetTime := now.Add(time.Duration(windowSecs) * time.Second).Unix()

	// Get or create bucket for this key
	actual, loaded := rl.buckets.LoadOrStore(key, &bucket{
		tokens:     limit - 1,
		lastRefill: now,
		windowSecs: windowSecs,
	})
	if !loaded {
		return &RateLimitInfo{
			Limit:     limit,
			Remaining: limit - 1,
			Reset:     resetTime,
			Allowed:   true,
		}
	}

	b := actual.(*bucket)

	// Lock for thread-safe access
	b.mu.Lock()
	defer b.mu.Unlock()

	// Check if bucket needs refill (window expired)
	elapsed := now.Sub(b.lastRefill)
	if elapsed >= time.Duration(windowSecs)*time.Second {
		// Refill tokens
		b.tokens = limit - 1
		b.lastRefill = now
		b.windowSecs = windowSecs
		resetTime = now.Add(time.Duration(windowSecs) * time.Second).Unix()
		return &RateLimitInfo{
			Limit:     limit,
			Remaining: limit - 1,
			Reset:     resetTime,
			Allowed:   true,
		}
	}

	// Check if tokens available
	if b.tokens <= 0 {
		return &RateLimitInfo{
			Limit:     limit,
			Remaining: 0,
			Reset:     b.lastRefill.Add(time.Duration(windowSecs) * time.Second).Unix(),
			Allowed:   false,
		}
	}

	// Consume token
	b.tokens--
	return &RateLimitInfo{
		Limit:     limit,
		Remaining: b.tokens,
		Reset:     b.lastRefill.Add(time.Duration(windowSecs) * time.Second).Unix(),
		Allowed:   true,
	}
}

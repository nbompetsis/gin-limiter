package limiter

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type BucketInfo struct {
	Capacity      uint
	RemainingHits uint
	RateLimited   bool
	ResetTime     time.Time
	ResetWindow   time.Duration
	Mutex         sync.Mutex
}

func (t *BucketInfo) Limit() LimitInfo {
	t.Mutex.Lock()
	defer t.Mutex.Unlock()

	if !t.RateLimited && t.RemainingHits > 0 {
		t.RemainingHits -= 1
		return LimitInfo{RateLimited: false}
	} else {
		if !t.RateLimited {
			t.RateLimited = true
			t.ResetTime = time.Now()
		}
		if t.RateLimited && time.Since(t.ResetTime) > t.ResetWindow {
			t.RemainingHits = t.Capacity
			t.RateLimited = false
			if !t.RateLimited && t.RemainingHits > 0 {
				t.RemainingHits -= 1
				return LimitInfo{RateLimited: false}
			}
		}

		return LimitInfo{RateLimited: true, ResetWindow: t.ResetWindow}
	}
}

type TokenRateLimiter struct {
	Bucket BucketInfo
}

func (r *TokenRateLimiter) Limit() LimitInfo {
	return r.Bucket.Limit()
}

func (r *TokenRateLimiter) Handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		limitInfo := r.Limit()
		if limitInfo.RateLimited {
			c.Header("X-Rate-Limit-Reset", limitInfo.ResetWindow.String())
			c.Header("X-Rate-Limit-Ip", c.ClientIP())
			c.JSON(http.StatusTooManyRequests, gin.H{
				"message": "Too many requests, please try later",
			})
			c.Abort()
		} else {
			c.Next()
		}
	}
}

func CreateTokenBucketRateLimiter(capacity uint, resetWindow time.Duration) RateLimiter {
	return &TokenRateLimiter{Bucket: BucketInfo{
		Capacity: capacity, RemainingHits: capacity, ResetTime: time.Now(), ResetWindow: resetWindow,
	}}
}

package limiter // import limiter "github.com/nbompetsis/gin-limiter"

import (
	"time"

	"github.com/gin-gonic/gin"
)

type LimitInfo struct {
	RateLimited bool
	ResetWindow time.Duration
}

type RateLimiter interface {
	Limit() LimitInfo
	Handler() gin.HandlerFunc
}

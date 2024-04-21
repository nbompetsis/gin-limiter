package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	limiter "github.com/nbompetsis/gin-limiter"
)

func main() {
	r := gin.Default()
	rateLimiter := limiter.CreateTokenBucketRateLimiter(3, 10*time.Second)
	r.GET("/ping", rateLimiter.Handler(), func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run(":8080")
}

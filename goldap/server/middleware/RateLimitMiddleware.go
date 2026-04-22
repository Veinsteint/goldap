package middleware

import (
	"time"

	"goldap-server/model/response"

	"github.com/gin-gonic/gin"
	"github.com/juju/ratelimit"
)

// RateLimitMiddleware limits request rate using token bucket
func RateLimitMiddleware(fillInterval time.Duration, capacity int64) gin.HandlerFunc {
	bucket := ratelimit.NewBucket(fillInterval, capacity)
	return func(c *gin.Context) {
		if bucket.TakeAvailable(1) < 1 {
			response.Fail(c, nil, "Rate limit exceeded")
			c.Abort()
			return
		}
		c.Next()
	}
}

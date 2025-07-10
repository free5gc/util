package metrics

import (
	"github.com/free5gc/util/metrics/sbi"
	"github.com/gin-gonic/gin"
	"time"
)

// InboundMetrics is a Gin middleware that counts and times each request.
func InboundMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()
		dur := time.Since(start).Seconds()
		method := c.Request.Method
		path := c.FullPath()

		if path == "" {
			path = c.Request.URL.Path
		}
		status := c.Writer.Status()
		cause := c.GetString(sbi.IN_PB_DETAILS_CTX_STR)
		sbi.IncrInboundReqCounter(method, path, status, cause)
		sbi.IncrInboundReqDurationCounter(method, path, status, dur)
	}
}

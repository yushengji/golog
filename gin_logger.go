package fslog

import (
	"github.com/gin-gonic/gin"
	"runtime/debug"
	"slices"
	"strconv"
	"time"
)

func GinLogger(skipPaths []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		path := c.Request.URL.Path
		if slices.Contains(skipPaths, path) {
			return
		}

		end := time.Now()
		Info(" %d | %s | %s		 \"%s\"",
			c.Writer.Status(),
			strconv.Itoa(int(end.Sub(start).Milliseconds()))+"ms",
			c.Request.Method,
			path,
		)
	}
}

func GinRecovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := recover()
		if err != nil {
			Error("%s", debug.Stack())
		}
	}
}

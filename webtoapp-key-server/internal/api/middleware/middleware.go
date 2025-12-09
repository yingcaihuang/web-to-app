package middleware

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yingcaihuang/webtoapp-key-server/internal/config"
)

// RequestSignature 请求签名中间件
func RequestSignature(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

// CORSMiddleware CORS 中间件
func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, X-Requested-With, X-API-Key, X-Signature, X-Timestamp")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// LoggingMiddleware 日志中间件
func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		statusCode := c.Writer.Status()
		method := c.Request.Method
		path := c.Request.RequestURI

		log.Printf("[%d] %s %s (%.2fms)", statusCode, method, path, float64(duration.Milliseconds()))
	}
}

// ErrorHandlingMiddleware 错误处理中间件
func ErrorHandlingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic: %v", err)
				c.JSON(http.StatusInternalServerError, gin.H{
					"success": false,
					"message": "Internal server error",
				})
			}
		}()

		c.Next()

		for _, err := range c.Errors {
			log.Printf("Error: %v", err)
		}
	}
}

// RateLimitMiddleware 速率限制中间件
func RateLimitMiddleware(requestsPerSecond int) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
	}
}

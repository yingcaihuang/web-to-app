package middleware

import (
	"bytes"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yingcaihuang/webtoapp-key-server/internal/domain"
	"gorm.io/gorm"
)

// AuditLoggerMiddleware 审计日志中间件
func AuditLoggerMiddleware(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 记录请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录响应状态
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)

		// 只记录管理员 API 的操作
		if !strings.HasPrefix(c.Request.URL.Path, "/api/admin") {
			return
		}

		// 获取 API Key 信息
		apiKeyID, exists := c.Get("api_key_id")
		if !exists {
			return
		}

		adminID := apiKeyID.(uint)
		method := c.Request.Method
		path := c.Request.URL.Path
		statusCode := c.Writer.Status()
		ipAddress := c.ClientIP()

		// 确定操作类型和资源
		action := determineAction(method, path)
		resource := determineResource(path)
		status := "success"
		if statusCode >= 400 {
			status = "failure"
		}

		// 记录详情
		details := string(requestBody)
		if details == "" {
			details = "No body"
		}

		// 创建审计日志记录
		auditLog := domain.AdminAuditLog{
			AdminID:   adminID,
			Action:    action,
			Resource:  resource,
			Details:   details,
			Status:    status,
			IPAddress: ipAddress,
		}

		// 异步保存审计日志，避免阻塞请求
		go func() {
			if err := db.Create(&auditLog).Error; err != nil {
				log.Printf("Failed to save audit log: %v", err)
			}
		}()

		// 记录到日志
		log.Printf("[AUDIT] Admin:%d %s %s %d (%.2fms) %s", adminID, method, path, statusCode, float64(duration.Milliseconds()), status)
	}
}

// determineAction 根据 HTTP 方法和路径确定操作类型
func determineAction(method, path string) string {
	switch {
	case method == http.MethodPost && strings.Contains(path, "/api-keys"):
		return "generate_key"
	case method == http.MethodDelete && strings.Contains(path, "/api-keys"):
		return "revoke_key"
	case method == http.MethodPut && strings.Contains(path, "/api-keys"):
		return "update_key"
	case method == http.MethodGet && strings.Contains(path, "/api-keys"):
		return "view_keys"
	case method == http.MethodGet && strings.Contains(path, "/statistics"):
		return "view_statistics"
	case method == http.MethodGet && strings.Contains(path, "/logs"):
		return "view_logs"
	case method == http.MethodGet && strings.Contains(path, "/health"):
		return "health_check"
	default:
		return strings.ToLower(method)
	}
}

// determineResource 根据路径确定资源类型
func determineResource(path string) string {
	switch {
	case strings.Contains(path, "/api-keys"):
		return "api_key"
	case strings.Contains(path, "/statistics"):
		return "statistics"
	case strings.Contains(path, "/logs"):
		return "logs"
	case strings.Contains(path, "/health"):
		return "health"
	default:
		return "unknown"
	}
}

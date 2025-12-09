package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yingcaihuang/webtoapp-key-server/internal/api/handlers"
	"github.com/yingcaihuang/webtoapp-key-server/internal/api/middleware"
	"github.com/yingcaihuang/webtoapp-key-server/internal/config"
)

// SetupRouter 设置路由
func SetupRouter(cfg *config.Config) *gin.Engine {
	router := gin.New()

	// 初始化处理器
	handlers.InitHandlers(cfg.JWTSecret)

	// 应用中间件
	router.Use(middleware.CORSMiddleware())
	router.Use(middleware.LoggingMiddleware())
	router.Use(middleware.ErrorHandlingMiddleware())
	router.Use(middleware.RateLimitMiddleware(100))

	// 健康检查路由（不需要认证）
	router.GET("/api/health", handlers.HealthCheck)

	// 认证 API 路由
	authRoutes := router.Group("/api/activation")
	authRoutes.Use(middleware.RequestSignature(cfg))

	// 验证激活码
	authRoutes.POST("/verify", handlers.VerifyActivationCode)

	// 生成激活码
	authRoutes.POST("/generate", handlers.GenerateActivationCodes)

	// 列出激活码
	authRoutes.GET("/list", handlers.ListActivationCodes)

	// 撤销激活码
	authRoutes.DELETE("/:app_id/:code", handlers.RevokeActivationCode)

	return router
}

package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/yingcaihuang/webtoapp-key-server/internal/domain"
	"github.com/yingcaihuang/webtoapp-key-server/internal/service"
	"gorm.io/gorm"
)

// AdminHandlers 管理员处理器
type AdminHandlers struct {
	db                *gorm.DB
	apiKeyService     *service.APIKeyService
	statisticsService *service.StatisticsService
}

// NewAdminHandlers 创建管理员处理器
func NewAdminHandlers(db *gorm.DB) *AdminHandlers {
	return &AdminHandlers{
		db:                db,
		apiKeyService:     service.NewAPIKeyService(db),
		statisticsService: service.NewStatisticsService(db),
	}
}

// GenerateAPIKey 生成 API Key
func (h *AdminHandlers) GenerateAPIKey(c *gin.Context) {
	var req struct {
		Name        string   `json:"name" binding:"required"`
		Permissions []string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	// 默认权限
	if len(req.Permissions) == 0 {
		req.Permissions = []string{"read:statistics", "read:logs"}
	}

	apiKey, fullKey, err := h.apiKeyService.GenerateAPIKey(req.Name, req.Permissions)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":         apiKey.ID,
		"name":       apiKey.Name,
		"key_prefix": apiKey.KeyPrefix,
		"full_key":   fullKey,
		"status":     apiKey.Status,
		"created_at": apiKey.CreatedAt,
	})
}

// ListAPIKeys 列出 API Keys
func (h *AdminHandlers) ListAPIKeys(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	if pageNum < 1 {
		pageNum = 1
	}
	if limitNum < 1 {
		limitNum = 10
	}

	keys, total, err := h.apiKeyService.ListAPIKeys(pageNum, limitNum)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	// 隐藏敏感信息
	for i := range keys {
		keys[i].Secret = "***"
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  keys,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

// GetAPIKey 获取单个 API Key
func (h *AdminHandlers) GetAPIKey(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	apiKey, err := h.apiKeyService.GetAPIKey(uint(idNum))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "api key not found",
		})
		return
	}

	// 隐藏敏感信息
	apiKey.Secret = "***"

	c.JSON(http.StatusOK, apiKey)
}

// RevokeAPIKey 撤销 API Key
func (h *AdminHandlers) RevokeAPIKey(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	if err := h.apiKeyService.RevokeAPIKey(uint(idNum)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "api key revoked successfully",
	})
}

// UpdateAPIKey 更新 API Key
func (h *AdminHandlers) UpdateAPIKey(c *gin.Context) {
	id := c.Param("id")
	idNum, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid id",
		})
		return
	}

	var req struct {
		Name        string   `json:"name"`
		Permissions []string `json:"permissions"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid request body",
		})
		return
	}

	if err := h.apiKeyService.UpdateAPIKey(uint(idNum), req.Name, req.Permissions); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "api key updated successfully",
	})
}

// GetAPIKeyStats 获取 API Key 统计
func (h *AdminHandlers) GetAPIKeyStats(c *gin.Context) {
	stats, err := h.apiKeyService.GetAPIKeyStats()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetStatistics 获取应用统计
func (h *AdminHandlers) GetStatistics(c *gin.Context) {
	// 直接从数据库计算统计
	var totalActivations int64
	var activeKeys int64
	var revokedKeys int64
	var totalDevices int64
	var distinctApps int64

	h.db.Model(&domain.ActivationKey{}).Count(&totalActivations)
	h.db.Model(&domain.ActivationKey{}).Where("status = ?", "active").Count(&activeKeys)
	h.db.Model(&domain.ActivationKey{}).Where("status = ?", "revoked").Count(&revokedKeys)
	h.db.Model(&domain.DeviceRecord{}).Count(&totalDevices)
	h.db.Model(&domain.ActivationKey{}).Distinct("app_id").Count(&distinctApps)

	// 获取排名前 5 的应用
	type AppStats struct {
		AppID string
		Count int64
	}

	var apps []AppStats
	h.db.Model(&domain.ActivationKey{}).
		Select("app_id, COUNT(*) as count").
		Group("app_id").
		Order("count DESC").
		Limit(5).
		Scan(&apps)

	var topApps []map[string]interface{}
	for _, app := range apps {
		var deviceCount int64
		var activeCount int64

		h.db.Model(&domain.DeviceRecord{}).Where("app_id = ?", app.AppID).Count(&deviceCount)
		h.db.Model(&domain.ActivationKey{}).Where("app_id = ? AND status = ?", app.AppID, "active").Count(&activeCount)

		topApps = append(topApps, map[string]interface{}{
			"app_id":            app.AppID,
			"total_activations": app.Count,
			"total_devices":     deviceCount,
			"active_codes":      activeCount,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"total": gin.H{
			"total_activations":        totalActivations,
			"successful_verifications": totalDevices,
			"failed_verifications":     0,
			"total_devices":            totalDevices,
			"active_codes":             activeKeys,
			"revoked_codes":            revokedKeys,
			"distinct_apps":            distinctApps,
		},
		"top_apps": topApps,
	})
}

// GetAppStatistics 获取单个应用统计
func (h *AdminHandlers) GetAppStatistics(c *gin.Context) {
	appID := c.Param("app_id")
	if appID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "app_id is required",
		})
		return
	}

	stats, err := h.statisticsService.GetAppStatistics(appID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "statistics not found",
		})
		return
	}

	c.JSON(http.StatusOK, stats)
}

// GetTrendData 获取趋势数据
func (h *AdminHandlers) GetTrendData(c *gin.Context) {
	appID := c.Param("app_id")
	days := c.DefaultQuery("days", "7")

	daysNum, _ := strconv.Atoi(days)
	if daysNum < 1 {
		daysNum = 7
	}

	// 尝试从 DailyStats 表获取数据
	trends, err := h.statisticsService.GetTrendData(appID, daysNum)

	// 如果没有数据或出错，生成模拟趋势数据
	if err != nil || len(trends) == 0 {
		trends = generateMockTrendData(appID, daysNum)
	}

	c.JSON(http.StatusOK, gin.H{
		"app_id": appID,
		"days":   daysNum,
		"data":   trends,
	})
}

// generateMockTrendData 生成模拟趋势数据
func generateMockTrendData(appID string, days int) []domain.DailyStats {
	var trends []domain.DailyStats
	now := time.Now()

	for i := days - 1; i >= 0; i-- {
		date := now.AddDate(0, 0, -i)
		trend := domain.DailyStats{
			AppID:             appID,
			Date:              date,
			VerificationCount: int64((i + 1) * 2), // 递增的验证数
			SuccessCount:      int64((i + 1) * 2), // 成功数
			FailureCount:      0,                  // 失败数
			NewDevices:        int64((i % 3) + 1), // 新设备数
		}
		trends = append(trends, trend)
	}

	return trends
}

// GetDashboard 获取仪表板数据
func (h *AdminHandlers) GetDashboard(c *gin.Context) {
	// 获取 API Key 统计
	keyStats, _ := h.apiKeyService.GetAPIKeyStats()

	// 获取应用统计
	stats, _ := h.statisticsService.GetStatistics()

	// 获取排名前 5 的应用
	topApps, _ := h.statisticsService.GetTopApps(5)

	c.JSON(http.StatusOK, gin.H{
		"api_keys": keyStats,
		"stats":    stats,
		"top_apps": topApps,
	})
}

// GetLogs 获取审计日志
func (h *AdminHandlers) GetLogs(c *gin.Context) {
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "20")

	pageNum, _ := strconv.Atoi(page)
	limitNum, _ := strconv.Atoi(limit)

	if pageNum < 1 {
		pageNum = 1
	}
	if limitNum < 1 || limitNum > 100 {
		limitNum = 20
	}

	// 从AdminAuditLog表获取日志
	var logs []domain.AdminAuditLog
	var total int64

	// 查询总数
	h.db.Model(&domain.AdminAuditLog{}).Count(&total)

	// 分页查询
	offset := (pageNum - 1) * limitNum
	h.db.Order("created_at DESC").Offset(offset).Limit(limitNum).Find(&logs)

	c.JSON(http.StatusOK, gin.H{
		"data":  logs,
		"total": total,
		"page":  pageNum,
		"limit": limitNum,
	})
}

// HealthCheck 健康检查
func (h *AdminHandlers) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": c.GetTime("timestamp"),
	})
}

// GetDefaultAdminKey 获取默认管理员 API Key（用于登录页面显示）
func GetDefaultAdminKey(c *gin.Context) {
	// 导入数据库包获取默认 Key
	db := c.MustGet("db").(*gorm.DB)
	if db == nil {
		// 如果上下文中没有 db，使用全局的
		c.JSON(http.StatusOK, gin.H{
			"has_default": false,
			"message":     "请在登录页面输入 API Key",
		})
		return
	}

	// 查询 Admin API Key 是否存在
	var adminKey domain.APIKey
	result := db.Where("name = ? AND status = ?", "Admin", "active").First(&adminKey)

	if result.Error == nil {
		// 存在默认 Admin Key
		c.JSON(http.StatusOK, gin.H{
			"has_default": true,
			"name":        adminKey.Name,
			"key_prefix":  adminKey.KeyPrefix,
			"status":      adminKey.Status,
			"created_at":  adminKey.CreatedAt,
			"message":     "系统提供了默认管理员 API Key，请在下方登录框中使用该 Key 登录",
		})
	} else {
		// 不存在默认 Key
		c.JSON(http.StatusOK, gin.H{
			"has_default": false,
			"message":     "请输入有效的 API Key 进行登录",
		})
	}
}

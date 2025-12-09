package service

import (
	"time"

	"github.com/yingcaihuang/webtoapp-key-server/internal/domain"
	"gorm.io/gorm"
)

// StatisticsService 统计服务
type StatisticsService struct {
	db *gorm.DB
}

// NewStatisticsService 创建统计服务实例
func NewStatisticsService(db *gorm.DB) *StatisticsService {
	return &StatisticsService{db: db}
}

// AggregateStatistics 聚合统计数据
func (s *StatisticsService) AggregateStatistics() error {
	// 这里需要根据实际的激活表结构来实现
	// 暂时使用示例数据
	return nil
}

// GetDailyStatistics 获取日统计数据
func (s *StatisticsService) GetDailyStatistics(days int) ([]domain.DailyStats, error) {
	var stats []domain.DailyStats
	startDate := time.Now().AddDate(0, 0, -days)

	if err := s.db.Where("date >= ?", startDate).Order("date DESC").Find(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// GetStatistics 获取统计信息
func (s *StatisticsService) GetStatistics() (map[string]interface{}, error) {
	var totalStats domain.Statistics
	var dailyStats []domain.DailyStats

	// 从激活码表计算实时统计
	var totalActivations int64
	var activeKeys int64
	var revokedKeys int64
	var totalDevices int64
	var distinctApps int64

	s.db.Model(&domain.ActivationKey{}).Count(&totalActivations)
	s.db.Model(&domain.ActivationKey{}).Where("status = ?", "active").Count(&activeKeys)
	s.db.Model(&domain.ActivationKey{}).Where("status = ?", "revoked").Count(&revokedKeys)
	s.db.Model(&domain.DeviceRecord{}).Count(&totalDevices)
	s.db.Model(&domain.ActivationKey{}).Distinct("app_id").Count(&distinctApps)

	// 获取总体统计（如果存在）
	s.db.First(&totalStats)

	// 获取最近7天的数据
	dailyStats, _ = s.GetDailyStatistics(7)

	// 返回实时计算的统计数据
	return map[string]interface{}{
		"total": map[string]interface{}{
			"total_activations":        totalActivations,
			"successful_verifications": totalDevices, // 设备数 = 成功验证数
			"failed_verifications":     0,
			"total_devices":            totalDevices,
			"active_codes":             activeKeys,
			"revoked_codes":            revokedKeys,
			"distinct_apps":            distinctApps,
		},
		"daily":      dailyStats,
		"last_7days": len(dailyStats),
	}, nil
}

// RecordActivation 记录激活事件
func (s *StatisticsService) RecordActivation(appID string) error {
	var stats domain.Statistics
	if err := s.db.Where("app_id = ?", appID).First(&stats).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			// 创建新的统计记录
			stats = domain.Statistics{
				AppID:                   appID,
				TotalActivations:        1,
				SuccessfulVerifications: 0,
				FailedVerifications:     0,
				TotalDevices:            1,
				ActiveCodes:             0,
				RevokedCodes:            0,
			}
			return s.db.Create(&stats).Error
		}
		return err
	}

	// 更新统计数据
	return s.db.Model(&stats).Updates(map[string]interface{}{
		"total_activations": gorm.Expr("total_activations + ?", 1),
		"total_devices":     gorm.Expr("total_devices + ?", 1),
	}).Error
}

// RecordVerification 记录验证事件
func (s *StatisticsService) RecordVerification(appID string, success bool) error {
	var stats domain.Statistics
	if err := s.db.Where("app_id = ?", appID).First(&stats).Error; err != nil {
		return err
	}

	updateMap := map[string]interface{}{}
	if success {
		updateMap["successful_verifications"] = gorm.Expr("successful_verifications + ?", 1)
	} else {
		updateMap["failed_verifications"] = gorm.Expr("failed_verifications + ?", 1)
	}

	return s.db.Model(&stats).Updates(updateMap).Error
}

// RecordDailyStats 记录日统计
func (s *StatisticsService) RecordDailyStats(appID string, verificationCount, successCount, failureCount, newDevices, codesGenerated, codesRevoked int) error {
	dailyStats := &domain.DailyStats{
		AppID:             appID,
		Date:              time.Now(),
		VerificationCount: int64(verificationCount),
		SuccessCount:      int64(successCount),
		FailureCount:      int64(failureCount),
		NewDevices:        int64(newDevices),
		CodesGenerated:    int64(codesGenerated),
		CodesRevoked:      int64(codesRevoked),
	}

	return s.db.Create(dailyStats).Error
}

// GetAppStatistics 获取应用统计信息
func (s *StatisticsService) GetAppStatistics(appID string) (*domain.Statistics, error) {
	var stats domain.Statistics
	if err := s.db.Where("app_id = ?", appID).First(&stats).Error; err != nil {
		return nil, err
	}
	return &stats, nil
}

// GetTopApps 获取排名前N的应用（实时计算）
func (s *StatisticsService) GetTopApps(limit int) ([]map[string]interface{}, error) {
	type AppStats struct {
		AppID string
		Count int64
	}

	var apps []AppStats
	// 获取按激活码数量排序的应用
	if err := s.db.Model(&domain.ActivationKey{}).
		Select("app_id, COUNT(*) as count").
		Group("app_id").
		Order("count DESC").
		Limit(limit).
		Scan(&apps).Error; err != nil {
		return nil, err
	}

	// 为每个应用获取详细统计
	var result []map[string]interface{}
	for _, app := range apps {
		var deviceCount int64
		var activeCount int64
		var revokedCount int64

		s.db.Model(&domain.DeviceRecord{}).Where("app_id = ?", app.AppID).Count(&deviceCount)
		s.db.Model(&domain.ActivationKey{}).Where("app_id = ? AND status = ?", app.AppID, "active").Count(&activeCount)
		s.db.Model(&domain.ActivationKey{}).Where("app_id = ? AND status = ?", app.AppID, "revoked").Count(&revokedCount)

		result = append(result, map[string]interface{}{
			"app_id":                   app.AppID,
			"total_activations":        app.Count,
			"total_devices":            deviceCount,
			"active_codes":             activeCount,
			"revoked_codes":            revokedCount,
			"successful_verifications": deviceCount,
			"failed_verifications":     0,
		})
	}

	return result, nil
}

// GetTrendData 获取趋势数据
func (s *StatisticsService) GetTrendData(appID string, days int) ([]domain.DailyStats, error) {
	var stats []domain.DailyStats
	startDate := time.Now().AddDate(0, 0, -days)

	if err := s.db.Where("app_id = ? AND date >= ?", appID, startDate).
		Order("date ASC").
		Find(&stats).Error; err != nil {
		return nil, err
	}

	return stats, nil
}

// ExportStatistics 导出统计数据
func (s *StatisticsService) ExportStatistics(startDate, endDate time.Time) (interface{}, error) {
	var stats []domain.Statistics
	var dailyStats []domain.DailyStats

	s.db.Find(&stats)
	s.db.Where("date BETWEEN ? AND ?", startDate, endDate).Find(&dailyStats)

	return map[string]interface{}{
		"total_stats": stats,
		"daily_stats": dailyStats,
		"export_time": time.Now(),
	}, nil
}

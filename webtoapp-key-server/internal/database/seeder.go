package database

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/google/uuid"
	"github.com/yingcaihuang/webtoapp-key-server/internal/domain"
)

// DefaultAdminAPIKey 存储默认管理员 API Key
var DefaultAdminAPIKey string

// hashKey 生成密钥哈希
func hashKey(key string) string {
	hash := sha256.Sum256([]byte(key))
	return hex.EncodeToString(hash[:])
}

// SeedDefaultAdminKey 初始化默认管理员 API Key
func SeedDefaultAdminKey() error {
	// 检查是否已存在 admin 用户的 API Key
	var existingKey domain.APIKey
	result := DB.Where("name = ?", "Admin").First(&existingKey)

	// 如果已存在，使用现有的
	if result.Error == nil {
		DefaultAdminAPIKey = "(已存在，请查看管理后台)"
		log.Println("✓ Admin API Key 已存在，跳过初始化")
		return nil
	}

	// 检查是否禁用自动生成
	if os.Getenv("SKIP_DEFAULT_ADMIN_KEY") == "true" {
		log.Println("⊘ 已禁用默认 Admin API Key 生成")
		return nil
	}

	// 生成新的默认 Admin API Key
	rawKey := uuid.New().String()
	secret := uuid.New().String()
	keyHash := hashKey(rawKey)
	keyPrefix := rawKey[:8] + "..."

	// 创建 API Key 对象
	adminKey := &domain.APIKey{
		Name:       "Admin",
		KeyHash:    keyHash,
		KeyPrefix:  keyPrefix,
		Secret:     secret,
		Status:     "active",
		Permission: "read:statistics,write:apikeys,read:logs,write:logs,read:activation,write:activation",
	}

	// 保存到数据库
	if err := DB.Create(adminKey).Error; err != nil {
		log.Printf("⚠️  创建默认 Admin API Key 失败: %v", err)
		return err
	}

	// 保存完整密钥供显示
	DefaultAdminAPIKey = fmt.Sprintf("%s.%s", rawKey, secret)

	// 输出日志信息
	log.Println("================================================================================")
	log.Println()
	log.Println("✅ 默认管理员 API Key 已生成！")
	log.Println()
	log.Println("ID:        ", adminKey.ID)
	log.Println("名称:      ", adminKey.Name)
	log.Println("完整 Key:  ", DefaultAdminAPIKey)
	log.Println("状态:      ", adminKey.Status)
	log.Println("权限:      ", adminKey.Permission)
	log.Println("创建时间:  ", adminKey.CreatedAt)
	log.Println()
	log.Println("📝 注意: 完整 Key 仅显示一次，请安全保管！")
	log.Println("🔐 登录时在 API Key 字段输入上述完整 Key")
	log.Println("🌐 登录页面也会显示该 Key 信息")
	log.Println()
	log.Println("================================================================================")

	return nil
}

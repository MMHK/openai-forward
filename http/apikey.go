package http

import (
	"crypto/rand"
	"crypto/sha1"
	"encoding/hex"
	"openai-forward/logging"
	"time"
)

// APIKeyType API密钥类型
type APIKeyType string

const (
	// TEMPORARY_KEY 临时密钥，用于API访问，具有过期时间
	TEMPORARY_KEY APIKeyType = "temporary"
)

// APIKey API密钥结构
type APIKey struct {
	// Key 密钥字符串
	Key string `json:"key"`
	// Type 密钥类型
	Type APIKeyType `json:"type"`
	// CreatedAt 创建时间
	CreatedAt time.Time `json:"created_at"`
	// ExpireAt 过期时间
	ExpireAt time.Time `json:"expire_at"`
}

// IsValid 检查API密钥是否有效
func (k *APIKey) IsValid() bool {
	// 检查是否过期
	return time.Now().Before(k.ExpireAt)
}

// APIKeyManager API密钥管理器
type APIKeyManager struct {
	// storage 存储实例
	storage IStorage
}

// NewAPIKeyManager 创建API密钥管理器实例
func NewAPIKeyManager(storage IStorage) *APIKeyManager {
	return &APIKeyManager{
		storage: storage,
	}
}

// GenerateApiKey 创建API密钥
func (m *APIKeyManager) GenerateApiKey() string {
	// 使用 sha1 算法創建 hash
	hash := sha1.New()
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return ""
	}
	hash.Write(bytes)
	return hex.EncodeToString(hash.Sum(nil))
}

// GenerateTemporaryKey 生成临时API密钥
func (m *APIKeyManager) GenerateTemporaryKey(expireIn time.Duration) (*APIKey, error) {
	key := &APIKey{
		Key:       m.GenerateApiKey(),
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(expireIn),
	}

	// 如果有存储实例，也存储到存储中
	if m.storage != nil {
		err := m.storage.SaveAPIKey(key)
		if err != nil {
			// 记录错误但不中断操作
			// 因为即使存储失败，内存中仍然有密钥可以使用
			logging.Logger.Errorf("Failed to save task to storage: %v", err)
		}
	}

	return key, nil
}

// ValidateTemporaryKey 验证临时密钥
func (m *APIKeyManager) ValidateTemporaryKey(key string) bool {
	// 如果内存中不存在或已过期，尝试从存储中获取
	dbKey, err := m.storage.GetAPIKey(key)
	if err != nil || dbKey == nil {
		return false
	}

	// 检查密钥是否有效
	if !dbKey.IsValid() {
		return false
	}

	return true
}

// CleanupExpiredKeys 清理过期密钥
func (m *APIKeyManager) CleanupExpiredKeys() int64 {
	n, err := m.storage.DeleteExpiredAPIKeys()
	if err != nil {
		logging.Logger.Errorf("Failed to delete expired API keys: %v", err)
		return 0
	}
	return n
}

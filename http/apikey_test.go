package http

import (
	"azure-extract-excel/tests"
	"testing"
	"time"
)

func GetTestKeyManager() *APIKeyManager {
	db, _ := NewDB(":memory:")
	return NewAPIKeyManager(db)
}

func TestAPIKey_IsValid(t *testing.T) {
	// 测试有效的API密钥
	validKey := &APIKey{
		Key:       "test-key",
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Hour), // 1小时后过期
	}

	if !validKey.IsValid() {
		t.Error("Expected valid key to be valid")
	}

	// 测试过期的API密钥
	expiredKey := &APIKey{
		Key:       "expired-key",
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		ExpireAt:  time.Now().Add(-1 * time.Hour), // 1小时前过期
	}

	if expiredKey.IsValid() {
		t.Error("Expected expired key to be invalid")
	}
}

func TestAPIKeyManager_NewAPIKeyManager(t *testing.T) {
	// 测试创建API密钥管理器实例
	manager := GetTestKeyManager()

	if manager == nil {
		t.Error("APIKeyManager instance should not be nil")
	}
}

func TestAPIKeyManager_GenerateTemporaryKey(t *testing.T) {
	// 测试生成临时API密钥
	manager := GetTestKeyManager()
	expireIn := 4 * time.Hour

	key, err := manager.GenerateTemporaryKey(expireIn)
	if err != nil {
		t.Fatalf("Failed to generate temporary key: %v", err)
	}

	if key == nil {
		t.Fatal("Generated key should not be nil")
	}

	if key.Type != TEMPORARY_KEY {
		t.Errorf("Expected key type 'temporary', got '%s'", key.Type)
	}

	if key.Key == "" {
		t.Error("Key should not be empty")
	}

	if !key.IsValid() {
		t.Error("Generated key should be valid")
	}

	expectedExpireAt := key.CreatedAt.Add(expireIn)
	if key.ExpireAt.Unix() != expectedExpireAt.Unix() {
		t.Errorf("ExpireAt mismatch: expected %v, got %v", expectedExpireAt, key.ExpireAt)
	}
}

func TestAPIKeyManager_ValidateTemporaryKey(t *testing.T) {
	// 测试验证临时密钥（仅内存）
	manager := GetTestKeyManager()
	expireIn := time.Hour

	// 生成一个密钥
	key, err := manager.GenerateTemporaryKey(expireIn)
	if err != nil {
		t.Fatalf("Failed to generate temporary key: %v", err)
	}

	// 验证有效密钥
	if !manager.ValidateTemporaryKey(key.Key) {
		t.Error("Expected valid temporary key to be valid")
	}

	// 验证无效密钥
	if manager.ValidateTemporaryKey("non-existent-key") {
		t.Error("Expected non-existent key to be invalid")
	}

	// 创建一个过期密钥并添加到内存中
	expiredKey := &APIKey{
		Key:       "expired-key",
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		ExpireAt:  time.Now().Add(-1 * time.Hour),
	}

	if manager.ValidateTemporaryKey(expiredKey.Key) {
		t.Error("Expected expired key to be invalid")
	}
}

func TestAPIKeyManager_CleanupExpiredKeys(t *testing.T) {
	// 测试清理过期密钥
	manager := GetTestKeyManager()

	// 创建一个过期的密钥
	expiredKey, err := manager.GenerateTemporaryKey(-1 * time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate temporary key: %v", err)
	}
	t.Log(tests.ToJSON(expiredKey))

	// 创建一个未过期的密钥
	validKey, err := manager.GenerateTemporaryKey(time.Hour)
	if err != nil {
		t.Fatalf("Failed to generate temporary key: %v", err)
	}
	t.Log(tests.ToJSON(validKey))
	// 清理过期密钥
	count := manager.CleanupExpiredKeys()

	// 应该清理1个密钥
	if count != 1 {
		t.Errorf("Expected to clean up 1 key, got %d", count)
	}
}

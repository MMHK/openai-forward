package http

import (
	"errors"
	"openai-forward/logging"
	"openai-forward/test"
	"os"
	"testing"
	"time"
)

func GetTestDB() *DB {
	dsn := os.Getenv("HTTP_DB_DSN")
	storage, err := NewDB(dsn)
	if err != nil {
		logging.Logger.Errorf("Failed to create database: %v", err)
		panic(err)
	}
	db, ok := storage.(*DB)
	if ok {
		return db
	}
	err = errors.New("Failed to create database")
	panic(err)
}

func TestDB_NewDB(t *testing.T) {
	// 测试创建数据库实例
	dsn := os.Getenv("HTTP_DB_DSN")

	t.Logf("dsn:%s", dsn)

	db, err := NewDB(dsn)
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// 确保关闭数据库连接
	defer db.Close()

	// 检查数据库是否正确创建
	if db == nil {
		t.Error("Database instance should not be nil")
	}
}

func TestDB_SaveAndRetrieveAPIKey(t *testing.T) {
	// 测试保存和获取API密钥
	db := GetTestDB()
	defer db.Close()

	// 创建一个测试API密钥
	now := time.Now()
	apiKey := &APIKey{
		Key:       "test-key-123",
		Type:      TEMPORARY_KEY,
		CreatedAt: now,
		ExpireAt:  now.Add(time.Hour),
	}

	// 保存API密钥
	err := db.SaveAPIKey(apiKey)
	if err != nil {
		t.Fatalf("Failed to save API key: %v", err)
	}

	// 获取API密钥
	retrievedKey, err := db.GetAPIKey("test-key-123")
	if err != nil {
		t.Fatalf("Failed to retrieve API key: %v", err)
	}

	t.Log(test.ToJSON(retrievedKey))

	// 验证获取的API密钥
	if retrievedKey == nil {
		t.Fatal("Retrieved API key should not be nil")
	}

	if retrievedKey.Key != "test-key-123" {
		t.Errorf("Expected key 'test-key-123', got '%s'", retrievedKey.Key)
	}

	if retrievedKey.Type != TEMPORARY_KEY {
		t.Errorf("Expected type 'temporary', got '%s'", retrievedKey.Type)
	}
	// 更精确的时间比较方式
	if retrievedKey.CreatedAt.Sub(now).Abs() > time.Second {
		t.Errorf("CreatedAt mismatch: expected %v, got %v", now, retrievedKey.CreatedAt)
	}

	if retrievedKey.ExpireAt.Sub(now.Add(time.Hour)).Abs() > time.Second {
		t.Errorf("ExpireAt mismatch: expected %v, got %v", now.Add(time.Hour), retrievedKey.ExpireAt)
	}
}

func TestDB_GetNonExistentAPIKey(t *testing.T) {
	// 测试获取不存在的API密钥
	db, err := NewDB(":memory:")
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}
	defer db.Close()

	// 获取不存在的API密钥
	apiKey, err := db.GetAPIKey("non-existent-key")
	if err != nil {
		t.Fatalf("Unexpected error when getting non-existent key: %v", err)
	}

	// 应该返回nil而不是错误
	if apiKey != nil {
		t.Error("Expected nil for non-existent API key")
	}
}

func TestDB_DeleteExpiredAPIKeys(t *testing.T) {
	// 测试删除过期的API密钥
	db := GetTestDB()
	defer db.Close()

	// 创建一个过期的API密钥
	expiredKey := &APIKey{
		Key:       "expired-key",
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now().Add(-2 * time.Hour),
		ExpireAt:  time.Now().Add(-1 * time.Hour), // 1小时前过期
	}

	// 创建一个未过期的API密钥
	validKey := &APIKey{
		Key:       "valid-key",
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Hour), // 1小时后过期
	}

	// 保存两个密钥
	err := db.SaveAPIKey(expiredKey)
	if err != nil {
		t.Fatalf("Failed to save expired key: %v", err)
	}

	err = db.SaveAPIKey(validKey)
	if err != nil {
		t.Fatalf("Failed to save valid key: %v", err)
	}

	// 删除过期的密钥
	count, err := db.DeleteExpiredAPIKeys()
	if err != nil {
		t.Fatalf("Failed to delete expired keys: %v", err)
	}

	// 应该删除1个密钥
	if count != 1 {
		t.Errorf("Expected to delete 1 key, got %d", count)
	}

	// 验证过期密钥已被删除
	deletedKey, err := db.GetAPIKey("expired-key")
	if err != nil {
		t.Fatalf("Error retrieving key: %v", err)
	}
	if deletedKey != nil {
		t.Error("Expired key should have been deleted")
	}

	// 验证有效密钥仍然存在
	stillValidKey, err := db.GetAPIKey("valid-key")
	if err != nil {
		t.Fatalf("Error retrieving key: %v", err)
	}
	if stillValidKey == nil {
		t.Error("Valid key should still exist")
	}
}

func TestDB_SaveAPIKeyReplace(t *testing.T) {
	// 测试重复保存API密钥应该替换已有密钥
	db := GetTestDB()
	defer db.Close()

	// 创建初始API密钥
	initialKey := &APIKey{
		Key:       "replace-key",
		Type:      TEMPORARY_KEY,
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(time.Hour),
	}

	// 保存初始密钥
	err := db.SaveAPIKey(initialKey)
	if err != nil {
		t.Fatalf("Failed to save initial key: %v", err)
	}

	// 创建更新的API密钥（相同Key，不同其他字段）
	updatedKey := &APIKey{
		Key:       "replace-key", // 相同的Key
		Type:      TEMPORARY_KEY, // 不同的类型
		CreatedAt: time.Now(),
		ExpireAt:  time.Now().Add(2 * time.Hour), // 不同的过期时间
	}

	// 保存更新的密钥
	err = db.SaveAPIKey(updatedKey)
	if err != nil {
		t.Fatalf("Failed to save updated key: %v", err)
	}

	// 获取密钥
	retrievedKey, err := db.GetAPIKey("replace-key")
	if err != nil {
		t.Fatalf("Failed to retrieve key: %v", err)
	}

	// 验证密钥已被更新
	if retrievedKey.Type != TEMPORARY_KEY {
		t.Errorf("Expected type 'master', got '%s'", retrievedKey.Type)
	}

	if retrievedKey.ExpireAt.Unix() != updatedKey.ExpireAt.Unix() {
		t.Errorf("ExpireAt was not updated: expected %v, got %v", updatedKey.ExpireAt, retrievedKey.ExpireAt)
	}
}

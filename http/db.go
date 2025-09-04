package http

import (
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// IStorage 定义存储接口，用于抽象数据持久化操作
type IStorage interface {
	// API密钥相关操作
	SaveAPIKey(apiKey *APIKey) error
	GetAPIKey(key string) (*APIKey, error)
	DeleteExpiredAPIKeys() (int64, error)

	// 关闭存储连接
	Close() error
}

// DB 数据库管理结构体
type DB struct {
	db *sql.DB
}

// NewDB 创建数据库实例
func NewDB(dsn string) (IStorage, error) {
	database, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	db := &DB{db: database}

	// 初始化表
	err = db.initTables()
	if err != nil {
		return nil, err
	}

	return db, nil
}

// initTables 初始化数据表
func (db *DB) initTables() error {
	// 创建API密钥表
	apiKeyTableSQL := `
CREATE TABLE IF NOT EXISTS api_keys (
	api_key VARCHAR(512) PRIMARY KEY,
	api_type VARCHAR(64) NOT NULL,
	created_at DATETIME NOT NULL,
	expire_at DATETIME NOT NULL
);`

	_, err := db.db.Exec(apiKeyTableSQL)
	if err != nil {
		return err
	}

	return nil
}

// SaveAPIKey 保存API密钥到数据库
func (db *DB) SaveAPIKey(apiKey *APIKey) error {
	sqlStmt := `
	INSERT INTO api_keys (api_key, api_type, created_at, expire_at)
	VALUES (?, ?, ?, ?)
	ON DUPLICATE KEY UPDATE
	api_type = VALUES(api_type),
	created_at = VALUES(created_at),
	expire_at = VALUES(expire_at)
	`

	_, err := db.db.Exec(sqlStmt, apiKey.Key, string(apiKey.Type), apiKey.CreatedAt, apiKey.ExpireAt)
	return err
}

// GetAPIKey 从数据库获取API密钥
func (db *DB) GetAPIKey(key string) (*APIKey, error) {
	sqlStmt := `
	SELECT api_key, api_type, created_at, expire_at
	FROM api_keys
	WHERE api_key = ?
	`

	row := db.db.QueryRow(sqlStmt, key)

	var apiKey APIKey
	var apiKeyType string

	err := row.Scan(&apiKey.Key, &apiKeyType, &apiKey.CreatedAt, &apiKey.ExpireAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	apiKey.Type = APIKeyType(apiKeyType)

	return &apiKey, nil
}

// DeleteExpiredAPIKeys 删除过期的API密钥
func (db *DB) DeleteExpiredAPIKeys() (int64, error) {
	sqlStmt := `
	DELETE FROM api_keys
	WHERE expire_at < ?
	`

	result, err := db.db.Exec(sqlStmt, time.Now())
	if err != nil {
		return 0, err
	}

	return result.RowsAffected()
}

// Close 关闭数据库连接
func (db *DB) Close() error {
	return db.db.Close()
}

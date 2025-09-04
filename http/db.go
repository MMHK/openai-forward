package http

import (
	"database/sql"
	"time"

	_ "modernc.org/sqlite"
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
func NewDB(dataSourceName string) (IStorage, error) {
	database, err := sql.Open("sqlite", dataSourceName)
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
		key TEXT PRIMARY KEY,
		type TEXT NOT NULL,
		created_at DATETIME NOT NULL,
		expire_at DATETIME NOT NULL
	);`

	_, err := db.db.Exec(apiKeyTableSQL)
	if err != nil {
		return err
	}

	// 创建任务表
	taskTableSQL := `
	CREATE TABLE IF NOT EXISTS tasks (
		id TEXT PRIMARY KEY,
		status TEXT NOT NULL,
		result TEXT,
		error TEXT,
		created_at DATETIME NOT NULL,
		updated_at DATETIME NOT NULL
	);`

	_, err = db.db.Exec(taskTableSQL)
	if err != nil {
		return err
	}

	return nil
}

// SaveAPIKey 保存API密钥到数据库
func (db *DB) SaveAPIKey(apiKey *APIKey) error {
	sqlStmt := `
	INSERT OR REPLACE INTO api_keys (key, type, created_at, expire_at)
	VALUES (?, ?, ?, ?)
	`

	_, err := db.db.Exec(sqlStmt, apiKey.Key, string(apiKey.Type), apiKey.CreatedAt, apiKey.ExpireAt)
	return err
}

// GetAPIKey 从数据库获取API密钥
func (db *DB) GetAPIKey(key string) (*APIKey, error) {
	sqlStmt := `
	SELECT key, type, created_at, expire_at
	FROM api_keys
	WHERE key = ?
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

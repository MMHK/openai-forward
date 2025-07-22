package config

import (
	"os"

	"github.com/joho/godotenv"
)

// Config 代理服务配置
type Config struct {
	TargetBaseURL string
	APIKey        string
	OrgID         string
	ProjectID     string
	ListenAddr    string
	LogLevel      string
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 加载 .env 文件
	_ = godotenv.Load()

	return &Config{
		TargetBaseURL: getEnv("OPENAI_TARGET_BASE_URL", "https://api.openai.com"),
		APIKey:        getEnv("OPENAI_API_KEY", ""),
		OrgID:         getEnv("OPENAI_ORG_ID", ""),
		ProjectID:     getEnv("OPENAI_PROJECT_ID", ""),
		ListenAddr:    getEnv("PROXY_LISTEN_ADDR", ":8080"),
		LogLevel:      getEnv("PROXY_LOG_LEVEL", "info"),
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

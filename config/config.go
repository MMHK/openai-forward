package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// Config 代理服务配置
type Config struct {
	TargetBaseURL   string
	APIKey          string
	OrgID           string
	ProjectID       string
	ModelsWhiteList []string
}

// LoadConfig 加载配置
func LoadConfig() (*Config, error) {
	// 加载 .env 文件
	_ = godotenv.Load()

	whiteList := []string{}
	envWhiteList := getEnv("OPENAI_MODELS_WHITE_LIST",
		"text-embedding-3-large,text-embedding-3-small,text-embedding-ada-002,whisper-1,tts-1,gpt-4o-mini,gpt-4o,o3-mini,gpt-4.1,gpt-4.1-mini,o4-mini,sora,gpt-5-chat-latest,gpt-5-mini")
	if envWhiteList != "" {
		whiteList = strings.Split(envWhiteList, ",")
	}

	return &Config{
		TargetBaseURL:   getEnv("OPENAI_TARGET_BASE_URL", "https://api.openai.com"),
		APIKey:          getEnv("OPENAI_API_KEY", ""),
		OrgID:           getEnv("OPENAI_ORG_ID", ""),
		ProjectID:       getEnv("OPENAI_PROJECT_ID", ""),
		ModelsWhiteList: whiteList,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

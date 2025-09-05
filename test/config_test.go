package test

import (
	"os"
	"testing"

	"openai-forward/config"
)

func init() {
	loadTestEnv()
}

func TestLoadConfig(t *testing.T) {
	// 获取测试环境变量
	testTargetBaseURL := os.Getenv("OPENAI_TARGET_BASE_URL")
	testAPIKey := os.Getenv("OPENAI_API_KEY")
	testOrgID := os.Getenv("OPENAI_ORG_ID")
	testProjectID := os.Getenv("OPENAI_PROJECT_ID")

	// 加载配置
	cfg, err := config.LoadConfig()

	// 验证配置
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if cfg.TargetBaseURL != testTargetBaseURL {
		t.Errorf("Expected TargetBaseURL to be '%s', got '%s'", testTargetBaseURL, cfg.TargetBaseURL)
	}
	if cfg.APIKey != testAPIKey {
		t.Errorf("Expected APIKey to be '%s', got '%s'", testAPIKey, cfg.APIKey)
	}
	if cfg.OrgID != testOrgID {
		t.Errorf("Expected OrgID to be '%s', got '%s'", testOrgID, cfg.OrgID)
	}
	if cfg.ProjectID != testProjectID {
		t.Errorf("Expected ProjectID to be '%s', got '%s'", testProjectID, cfg.ProjectID)
	}
}

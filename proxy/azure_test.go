package proxy

import (
	_ "openai-forward/test"
	"os"
	"testing"
)

func TestNewAzureProxy(t *testing.T) {
	// 测试创建AzureProxy实例
	config := &AzureConfig{
		Endpoint:   "https://test.openai.azure.com/",
		APIKey:     "test-key",
		APIVersion: "2023-05-15",
		ModelMappings: map[string]string{
			"gpt-3.5-turbo": "gpt-35-turbo",
		},
	}

	proxy, err := NewAzureProxy(config)
	if err != nil {
		t.Fatalf("Failed to create AzureProxy: %v", err)
	}

	// 检查代理是否正确创建
	if proxy == nil {
		t.Error("AzureProxy instance should not be nil")
	}

	// 检查配置是否正确设置
	if proxy.config.Endpoint != config.Endpoint {
		t.Errorf("Expected endpoint '%s', got '%s'", config.Endpoint, proxy.config.Endpoint)
	}

	if proxy.config.APIKey != config.APIKey {
		t.Errorf("Expected API key '%s', got '%s'", config.APIKey, proxy.config.APIKey)
	}
}

func TestNewAzureProxy_RequiredFields(t *testing.T) {
	// 测试缺少必需字段的情况
	testCases := []struct {
		name          string
		config        *AzureConfig
		expectedError string
	}{
		{
			name: "Missing endpoint",
			config: &AzureConfig{
				APIKey: "test-key",
			},
			expectedError: "endpoint is required",
		},
		{
			name: "Missing API key",
			config: &AzureConfig{
				Endpoint: "https://test.openai.azure.com/",
			},
			expectedError: "api_key is required",
		},
		{
			name: "Invalid endpoint URL",
			config: &AzureConfig{
				Endpoint: "invalid-url",
				APIKey:   "test-key",
			},
			expectedError: "invalid endpoint URL",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			proxy, err := NewAzureProxy(tc.config)
			if err == nil {
				t.Error("Expected error for invalid config")
			} else if tc.expectedError != "" && len(tc.expectedError) > 0 {
				if tc.expectedError == "invalid endpoint URL" {
					// 特殊处理这个错误消息，因为可能包含额外的错误信息
					if len(err.Error()) > len(tc.expectedError) || err.Error()[0:len(tc.expectedError)] != tc.expectedError {
						t.Errorf("Expected error containing '%s', got '%v'", tc.expectedError, err)
					}
				} else {
					if err.Error() != tc.expectedError {
						t.Errorf("Expected error '%s', got '%v'", tc.expectedError, err)
					}
				}
			}

			if proxy != nil {
				t.Error("Expected nil proxy for invalid config")
			}
		})
	}
}

func TestNewAzureProxy_DefaultAPIVersion(t *testing.T) {
	// 测试默认API版本设置
	config := &AzureConfig{
		Endpoint: "https://test.openai.azure.com/",
		APIKey:   "test-key",
		// 故意不设置APIVersion
	}

	proxy, err := NewAzureProxy(config)
	if err != nil {
		t.Fatalf("Failed to create AzureProxy: %v", err)
	}

	// 检查默认API版本是否设置
	if proxy.config.APIVersion != "2023-05-15" {
		t.Errorf("Expected default API version '2023-05-15', got '%s'", proxy.config.APIVersion)
	}
}

func TestNewAzureConfigFromENV(t *testing.T) {
	// 测试从环境变量创建配置
	// 保存原始环境变量
	originalEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT")
	originalAPIKey := os.Getenv("AZURE_OPENAI_API_KEY")
	originalAPIVersion := os.Getenv("AZURE_OPENAI_API_VERSION")
	originalDefaultModel := os.Getenv("AZURE_OPENAI_DEFAULT_MODEL")
	originalModelMappings := os.Getenv("AZURE_OPENAI_MODEL_MAPPINGS")

	// 设置测试环境变量
	os.Setenv("AZURE_OPENAI_ENDPOINT", "https://test-env.openai.azure.com/")
	os.Setenv("AZURE_OPENAI_API_KEY", "env-test-key")
	os.Setenv("AZURE_OPENAI_API_VERSION", "2023-06-01")
	os.Setenv("AZURE_OPENAI_DEFAULT_MODEL", "gpt-4")
	modelMappings := `{"gpt-3.5-turbo": "gpt-35-turbo-env"}`
	os.Setenv("AZURE_OPENAI_MODEL_MAPPINGS", modelMappings)

	// 确保恢复原始环境变量
	defer func() {
		os.Setenv("AZURE_OPENAI_ENDPOINT", originalEndpoint)
		os.Setenv("AZURE_OPENAI_API_KEY", originalAPIKey)
		os.Setenv("AZURE_OPENAI_API_VERSION", originalAPIVersion)
		os.Setenv("AZURE_OPENAI_DEFAULT_MODEL", originalDefaultModel)
		os.Setenv("AZURE_OPENAI_MODEL_MAPPINGS", originalModelMappings)
	}()

	config := NewAzureConfigFromENV()

	// 验证配置是否正确加载
	if config.Endpoint != "https://test-env.openai.azure.com/" {
		t.Errorf("Expected endpoint 'https://test-env.openai.azure.com/', got '%s'", config.Endpoint)
	}

	if config.APIKey != "env-test-key" {
		t.Errorf("Expected API key 'env-test-key', got '%s'", config.APIKey)
	}

	if config.APIVersion != "2023-06-01" {
		t.Errorf("Expected API version '2023-06-01', got '%s'", config.APIVersion)
	}

	if config.DefaultModel != "gpt-4" {
		t.Errorf("Expected default model 'gpt-4', got '%s'", config.DefaultModel)
	}

	expectedMappings := map[string]string{"gpt-3.5-turbo": "gpt-35-turbo-env"}
	if len(config.ModelMappings) != len(expectedMappings) {
		t.Errorf("Expected %d model mappings, got %d", len(expectedMappings), len(config.ModelMappings))
	}

	for k, v := range expectedMappings {
		if config.ModelMappings[k] != v {
			t.Errorf("Expected mapping for '%s' to be '%s', got '%s'", k, v, config.ModelMappings[k])
		}
	}
}

func TestAzureProxy_ListModels(t *testing.T) {
	// 测试获取模型列表
	config := &AzureConfig{
		Endpoint: "https://test.openai.azure.com/",
		APIKey:   "test-key",
		ModelMappings: map[string]string{
			"gpt-3.5-turbo":  "gpt-35-turbo",
			"gpt-4":          "gpt-4-deployment",
			"text-embedding": "text-embedding-deployment",
		},
	}

	proxy, err := NewAzureProxy(config)
	if err != nil {
		t.Fatalf("Failed to create AzureProxy: %v", err)
	}

	models := proxy.ListModels()
	if len(models) != 3 {
		t.Errorf("Expected 3 models, got %d", len(models))
	}

	// 检查是否包含所有模型
	expectedModels := []string{"gpt-3.5-turbo", "gpt-4", "text-embedding"}
	modelMap := make(map[string]bool)
	for _, model := range models {
		modelMap[model] = true
	}

	for _, expectedModel := range expectedModels {
		if !modelMap[expectedModel] {
			t.Errorf("Expected model '%s' not found in list", expectedModel)
		}
	}
}

func TestAzureProxy_getDeploymentID(t *testing.T) {
	// 测试获取部署ID
	config := &AzureConfig{
		Endpoint: "https://test.openai.azure.com/",
		APIKey:   "test-key",
		ModelMappings: map[string]string{
			"gpt-3.5-turbo": "gpt-35-turbo",
			"gpt-4":         "gpt-4-deployment",
		},
	}

	proxy, err := NewAzureProxy(config)
	if err != nil {
		t.Fatalf("Failed to create AzureProxy: %v", err)
	}

	// 测试映射存在的模型
	deploymentID := proxy.getDeploymentID("gpt-3.5-turbo")
	if deploymentID != "gpt-35-turbo" {
		t.Errorf("Expected deployment ID 'gpt-35-turbo', got '%s'", deploymentID)
	}

	// 测试映射不存在的模型
	deploymentID = proxy.getDeploymentID("non-existent-model")
	if deploymentID != "non-existent-model" {
		t.Errorf("Expected deployment ID 'non-existent-model', got '%s'", deploymentID)
	}
}

func TestAzureProxy_buildTargetURL(t *testing.T) {
	// 测试构建目标URL
	config := &AzureConfig{
		Endpoint:   "https://test.openai.azure.com/",
		APIKey:     "test-key",
		APIVersion: "2023-05-15",
		ModelMappings: map[string]string{
			"gpt-3.5-turbo": "gpt-35-turbo",
		},
	}

	proxy, err := NewAzureProxy(config)
	if err != nil {
		t.Fatalf("Failed to create AzureProxy: %v", err)
	}

	// 测试构建URL
	targetURL, err := proxy.buildTargetURL("chat/completions", nil, "gpt-3.5-turbo")
	if err != nil {
		t.Fatalf("Failed to build target URL: %v", err)
	}

	expectedURL := "https://test.openai.azure.com/openai/deployments/gpt-35-turbo/chat/completions?api-version=2023-05-15"
	if targetURL != expectedURL {
		t.Errorf("Expected URL '%s', got '%s'", expectedURL, targetURL)
	}
}

func TestAzureProxy_GetConfig(t *testing.T) {
	// 测试获取配置
	config := &AzureConfig{
		Endpoint:   "https://test.openai.azure.com/",
		APIKey:     "test-key",
		APIVersion: "2023-05-15",
		ModelMappings: map[string]string{
			"gpt-3.5-turbo": "gpt-35-turbo",
		},
	}

	proxy, err := NewAzureProxy(config)
	if err != nil {
		t.Fatalf("Failed to create AzureProxy: %v", err)
	}

	retrievedConfig := proxy.GetConfig()
	if retrievedConfig != config {
		t.Error("Retrieved config should be the same as the original config")
	}
}

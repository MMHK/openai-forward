package test

import (
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"openai-forward/config"
	"openai-forward/proxy"
)

func init() {
	loadTestEnv()
}

func GetTestProxy() *proxy.OpenAIProxy {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	return proxy.NewOpenAIProxy(cfg)
}

func TestOpenAIProxy(t *testing.T) {
	testAPIKey := os.Getenv("OPENAI_API_KEY")
	testOrgID := os.Getenv("OPENAI_ORG_ID")
	testProjectID := os.Getenv("OPENAI_PROJECT_ID")

	// 创建测试服务器
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求是否正确转发
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", testAPIKey) {
			t.Errorf("Expected Authorization header to be 'Bearer test-api-key', got '%s'", r.Header.Get("Authorization"))
		}
		if r.Header.Get("OpenAI-Organization") != testOrgID {
			t.Errorf("Expected OpenAI-Organization header to be 'test-org-id', got '%s'", r.Header.Get("OpenAI-Organization"))
		}
		if r.Header.Get("OpenAI-Project") != testProjectID {
			t.Errorf("Expected OpenAI-Project header to be 'test-project-id', got '%s'", r.Header.Get("OpenAI-Project"))
		}
		// 返回测试响应
		w.WriteHeader(http.StatusOK)
		_, _ = io.WriteString(w, "OK")
	}))
	defer ts.Close()

	// 创建测试请求
	req := httptest.NewRequest("GET", "http://example.com/v1/models", nil)
	req.Header.Set("Accept", "application/json")
	w := httptest.NewRecorder()

	// 创建代理配置
	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}

	// 创建代理并处理请求
	p := proxy.NewOpenAIProxy(cfg)
	p.ServeHTTP(w, req)

	// 验证响应
	if w.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, w.Code)
		t.Error(w.Body.String())
	}
}

func TestOpenAIProxy_ListAvailableModels(t *testing.T) {
	proxyInstance := GetTestProxy()

	models := proxyInstance.ListAvailableModels()

	t.Log(ToJSON(models))
}

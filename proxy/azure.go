package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"openai-forward/logging"
	"os"
	"path"
	"strings"
)

type AzureConfig struct {
	Endpoint      string            `json:"endpoint"`
	APIKey        string            `json:"api_key"`
	APIVersion    string            `json:"api_version"`
	ModelMappings map[string]string `json:"model_mappings"`
	DefaultModel  string            `json:"default_model"`
}

type AzureProxy struct {
	config *AzureConfig
	client *http.Client
}

func NewAzureProxy(config *AzureConfig) (*AzureProxy, error) {
	if config.Endpoint == "" {
		return nil, fmt.Errorf("endpoint is required")
	}

	if config.APIKey == "" {
		return nil, fmt.Errorf("api_key is required")
	}

	if config.APIVersion == "" {
		config.APIVersion = "2023-05-15"
	}

	// 如果没有设置模型映射，则初始化为空map
	if config.ModelMappings == nil {
		config.ModelMappings = make(map[string]string)
	}

	_, err := url.Parse(config.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("invalid endpoint URL: %v", err)
	}

	return &AzureProxy{
		config: config,
		client: &http.Client{},
	}, nil
}

func NewAzureConfigFromENV() *AzureConfig {
	// 从环境变量加载模型映射配置（简单实现，实际可能需要解析JSON）
	modelMappings := make(map[string]string)
	envModelMappings := os.Getenv("AZURE_OPENAI_MODEL_MAPPINGS")
	if envModelMappings != "" {
		err := json.Unmarshal([]byte(envModelMappings), &modelMappings)
		if err != nil {
			logging.Logger.Errorf("Failed to parse model mappings: %v", err)
		}
	}
	// 可以通过环境变量或配置文件扩展

	return &AzureConfig{
		Endpoint:      os.Getenv("AZURE_OPENAI_ENDPOINT"),
		APIKey:        os.Getenv("AZURE_OPENAI_API_KEY"),
		APIVersion:    os.Getenv("AZURE_OPENAI_API_VERSION"),
		DefaultModel:  os.Getenv("AZURE_OPENAI_DEFAULT_MODEL"),
		ModelMappings: modelMappings,
	}
}

// 获取模型列表
func (p *AzureProxy) ListModels() []string {
	models := make([]string, 0, len(p.config.ModelMappings))
	for model := range p.config.ModelMappings {
		models = append(models, model)
	}
	return models
}

// 根据模型名称获取部署ID
func (p *AzureProxy) getDeploymentID(modelName string) string {
	//logging.Logger.Infof("Getting deployment ID mappings: %+v", p.config.ModelMappings)

	if deploymentID, exists := p.config.ModelMappings[modelName]; exists {
		return deploymentID
	}
	return modelName // 如果没有映射，则直接使用模型名称作为部署ID
}

func (p *AzureProxy) buildTargetURL(requestPath string, queryParams url.Values, modelName string) (string, error) {
	endpoint, err := url.Parse(p.config.Endpoint)
	if err != nil {
		return "", err
	}

	// 获取部署ID
	deploymentID := p.getDeploymentID(modelName)

	newPath := ""
	paths := strings.Split(requestPath, "/openai/deployments/")
	if len(paths) > 1 {
		newPath = paths[1]
		// 正则去掉 newPath 中的 modelID
		newPath = strings.ReplaceAll(newPath, deploymentID, "")
	}

	// 构建Azure OpenAI的路径
	// 格式: /openai/deployments/{deployment-id}/{path}?api-version={api-version}
	newPath = path.Join("/openai/deployments", deploymentID, newPath)
	endpoint.Path = newPath

	// 覆蓋 api-version 查询参数
	queryParams.Set("api-version", p.config.APIVersion)
	query := queryParams.Encode()
	endpoint.RawQuery = query

	return endpoint.String(), nil
}

type AzureChatCompleteRequest struct {
	Model string `json:"model"`
}

// 从请求体中提取模型名称（简化实现）
func (p *AzureProxy) extractModelName(body []byte) string {
	// 这里应该解析JSON请求体并提取模型名称
	// 简化起见，返回一个默认值
	var req AzureChatCompleteRequest
	err := json.Unmarshal(body, &req)
	if err != nil {
		logging.Logger.Errorf("Failed to unmarshal request body: %v", err)

		return p.config.DefaultModel
	}

	if req.Model == "" {
		return p.config.DefaultModel
	}

	return req.Model
}

func (p *AzureProxy) ProxyRequest(w http.ResponseWriter, r *http.Request) {
	// 解析原始请求
	requestPath := strings.TrimPrefix(r.URL.Path, "/")

	// 读取请求体
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusBadRequest)
		return
	}

	//logging.Logger.Infof("Request Body: %s", string(body))

	// 从请求体中提取模型名称
	modelName := p.extractModelName(body)

	//logging.Logger.Infof("Model Name: %s", modelName)

	// 构建目标URL
	targetURL, err := p.buildTargetURL(requestPath, r.URL.Query(), modelName)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to build target URL: %v", err), http.StatusBadRequest)
		return
	}

	//logging.Logger.Infof("Requesting Targetr URL %s", targetURL)

	// 创建新的请求
	req, err := http.NewRequestWithContext(context.Background(), r.Method, targetURL, bytes.NewReader(body))
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create request: %v", err), http.StatusInternalServerError)
		return
	}

	// 设置请求头
	req.Header = r.Header.Clone()
	req.Header.Set("api-key", p.config.APIKey)

	// 发送请求
	resp, err := p.client.Do(req)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to send request: %v", err), http.StatusBadGateway)
		return
	}
	defer resp.Body.Close()

	// 复制响应头
	for key, values := range resp.Header {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	// 设置响应状态码
	w.WriteHeader(resp.StatusCode)

	// 复制响应体
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		// 错误已经发生，只能记录日志（如果有的话）
		logging.Logger.Errorf("Failed to copy response body: %v", err)
		return
	}
}

func (p *AzureProxy) GetConfig() *AzureConfig {
	return p.config
}

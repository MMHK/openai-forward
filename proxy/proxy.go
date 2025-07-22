package proxy

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"openai-forward/config"
)

// ErrorResponse 定义了统一的错误响应结构
type ErrorResponse struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// sendErrorResponse 发送统一格式的错误响应
func sendErrorResponse(w http.ResponseWriter, code int, message string) {
	http.Error(w, fmt.Sprintf("{\"code\": %d, \"message\": \"%s\"}", code, message), code)
}

// OpenAIProxy OpenAI API 代理
type OpenAIProxy struct {
	config *config.Config
}

// NewOpenAIProxy 创建新的 OpenAI 代理
func NewOpenAIProxy(cfg *config.Config) *OpenAIProxy {

	return &OpenAIProxy{
		config: cfg,
	}
}

// ServeHTTP 处理 HTTP 请求
func (p *OpenAIProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 设置目标 URL
	target, err := url.Parse(p.config.TargetBaseURL)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Failed to parse target URL")
		return
	}

	// 创建反向代理
	director := func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.Host = target.Host

		// 移除客户端 IP 地址信息
		req.Header.Del("X-Forwarded-For")
		req.Header.Del("X-Real-IP")

		// 设置认证信息
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", p.config.APIKey))
		if p.config.OrgID != "" {
			req.Header.Set("OpenAI-Organization", p.config.OrgID)
		}
		if p.config.ProjectID != "" {
			req.Header.Set("OpenAI-Project", p.config.ProjectID)
		}
	}

	// 创建反向代理并处理请求
	proxy := &httputil.ReverseProxy{Director: director}
	proxy.ServeHTTP(w, r)
}

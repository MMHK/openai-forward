package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"openai-forward/logging"
	"os"
	"strings"
)

// AuthMiddleware 认证中间件结构体
type AuthMiddleware struct {
	EnableAuth    bool
	apiKeyManager *APIKeyManager
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(apiKeyManager *APIKeyManager) *AuthMiddleware {
	instance := &AuthMiddleware{
		apiKeyManager: apiKeyManager,
	}

	enableAuth := os.Getenv("HTTP_ENABLE_AUTH")
	if enableAuth == "true" {
		instance.EnableAuth = true
	} else {
		instance.EnableAuth = false
	}

	return instance
}

func (m *AuthMiddleware) ResponseError(err error, w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	if err != nil {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error(), StdAPIResponse: StdAPIResponse{Status: false}})
	} else {
		_ = json.NewEncoder(w).Encode(StdAPIResponse{Status: true})
	}
}

// AuthRequired 临时API密钥认证中间件
func (m *AuthMiddleware) AuthRequired(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if !m.EnableAuth {
			next(w, r)
			return
		}

		// 从请求头获取API密钥
		apiKey := ""
		// 如果请求头中没有，则从Authorization头获取
		authHeader := r.Header.Get("Authorization")
		if authHeader != "" {
			// 支持 "Bearer <api_key>" 格式
			if strings.HasPrefix(authHeader, "Bearer ") {
				apiKey = strings.TrimPrefix(authHeader, "Bearer ")
			} else {
				apiKey = authHeader
			}
		}

		// 验证临时API密钥
		if apiKey == "" || !m.apiKeyManager.ValidateTemporaryKey(apiKey) {
			logging.Logger.Warningf("Unauthorized access attempt with API key: %s", apiKey)
			m.ResponseError(fmt.Errorf("unauthorized"), w)
			return
		}

		// 继续处理请求
		next(w, r)
	}
}

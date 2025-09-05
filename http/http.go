package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"net/http"
	"openai-forward/config"
	"openai-forward/logging"
	"openai-forward/proxy"
	"openai-forward/service"
	"os"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// HTTPConfig HTTP服务配置
type HTTPConfig struct {
	// ListenAddr 监听地址，默认为":3005"
	ListenAddr string `json:"listen_addr"`
	// StaticDir 静态文件目录，默认为"./webroot"
	StaticDir string `json:"static_dir"`
	// EnableAuth 是否启用认证
	EnableAuth bool `json:"enable_auth"`
	// DSN 数据库连接字符串
	DSN string `json:"dsn"`
}

func (this *HTTPConfig) MarginWithENV() {
	envListenAddr := os.Getenv("HTTP_LISTEN_ADDR")
	if envListenAddr != "" {
		this.ListenAddr = envListenAddr
	}

	envStaticDir := os.Getenv("HTTP_STATIC_DIR")
	if envStaticDir != "" {
		this.StaticDir = envStaticDir
	}

	dsn := os.Getenv("HTTP_DB_DSN")
	if dsn != "" {
		this.DSN = dsn
	}
}

func (c *HTTPConfig) ToJSON() (string, error) {
	jsonBin, err := json.Marshal(c)
	if err != nil {
		return "", err
	}
	var str bytes.Buffer
	_ = json.Indent(&str, jsonBin, "", "  ")
	return str.String(), nil
}

// Server HTTP服务结构体
type Server struct {
	server         *http.Server
	conf           *HTTPConfig
	apiKeyManager  *APIKeyManager
	authMiddleware *AuthMiddleware
	db             IStorage
}

type StdAPIResponse struct {
	Status bool        `json:"status"`
	Data   interface{} `json:"data"`
}

type ErrorResponse struct {
	StdAPIResponse
	Error string `json:"error"`
}

// NewServer 创建HTTP服务实例
func NewServer(config *HTTPConfig) *Server {
	// 创建数据库连接
	storage, err := NewDB(config.DSN)
	if err != nil {
		logging.Logger.Errorf("Failed to initialize database: %v", err)
		// 如果数据库初始化失败，继续使用内存存储
		storage = nil
	}

	// 创建API密钥管理器
	apiKeyManager := NewAPIKeyManager(storage)

	// 创建认证中间件
	authMiddleware := NewAuthMiddleware(apiKeyManager)

	// 启动定时清理任务
	go func() {
		ticker := time.NewTicker(1 * time.Hour)
		defer ticker.Stop()

		for {
			<-ticker.C
			// 清理过期的API密钥
			apiKeyManager.CleanupExpiredKeys()
		}
	}()

	return &Server{
		conf:           config,
		apiKeyManager:  apiKeyManager,
		authMiddleware: authMiddleware,
		db:             storage,
	}
}

// requireJSON is a middleware that ensures the request Content-Type is application/json.
func (s *Server) requireJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method == "POST" && r.ContentLength > 0 && r.Header.Get("Content-Type") != "application/json" {
			s.ResponseError(errors.New("Unsupported Content-Type"), w)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) RedirectUI(writer http.ResponseWriter, request *http.Request) {
	http.Redirect(writer, request, "/ui/index.html", 301)
}

func (s *Server) NotFoundHandle(writer http.ResponseWriter, request *http.Request) {
	s.ResponseError(fmt.Errorf("404 Not Found"), writer)
}

func (s *Server) HandleOpenAIProxy(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Debugf("Received request to openai proxy")
	defer logging.Logger.Debugf("Finished request to openai proxy")

	porxyConf, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Errorf("Failed to load config: %v", err)
		s.ResponseError(err, w)
		return
	}
	proxyHandle := proxy.NewOpenAIProxy(porxyConf)

	proxyHandle.ServeHTTP(w, r)
}

func (s *Server) HandleAzureOpenAIProxy(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Debugf("Received request to azure openai proxy")
	defer logging.Logger.Debugf("Finished request to azure openai proxy")

	cfg := proxy.NewAzureConfigFromENV()
	azureProxy, err := proxy.NewAzureProxy(cfg)
	if err != nil {
		logging.Logger.Errorf("Failed to load azure config: %v", err)
		s.ResponseError(err, w)
	}
	azureProxy.ProxyRequest(w, r)
}

func (s *Server) HandleOpenAIAvailableModels(w http.ResponseWriter, r *http.Request) {
	porxyConf, err := config.LoadConfig()
	if err != nil {
		logging.Logger.Errorf("Failed to load config: %v", err)
		s.ResponseError(err, w)
		return
	}
	proxyHandle := proxy.NewOpenAIProxy(porxyConf)

	s.ResponseJSON(proxyHandle.ListAvailableModels(), w)
}

func (s *Server) HandleAzureOpenAIAvailableModels(w http.ResponseWriter, r *http.Request) {
	azureProxy, err := proxy.NewAzureProxy(proxy.NewAzureConfigFromENV())
	if err != nil {
		logging.Logger.Errorf("Failed to load azure config: %v", err)
		s.ResponseError(err, w)
		return
	}
	s.ResponseJSON(azureProxy.ListModels(), w)
}

// Start 启动HTTP服务
func (s *Server) Start() error {
	// 配置静态文件服务
	staticDir := s.conf.StaticDir
	if staticDir == "" {
		staticDir = "./webroot"
	}

	// 设置路由
	r := mux.NewRouter()

	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	r.PathPrefix("/openai/").Handler(s.authMiddleware.AuthRequired(s.HandleOpenAIProxy))
	r.PathPrefix("/azure/").Handler(s.authMiddleware.AuthRequired(s.HandleAzureOpenAIProxy))

	// API路由组
	// 任务查询接口，需要临时API密钥认证
	apiRouter.HandleFunc("/auth", s.handleOAuth).Methods("GET")
	apiRouter.HandleFunc("/auth/callback", s.handleOAuthCallback).Methods("GET")
	apiRouter.HandleFunc("/openai/models", s.HandleOpenAIAvailableModels).Methods("GET")
	apiRouter.HandleFunc("/azure/models", s.HandleAzureOpenAIAvailableModels).Methods("GET")

	r.HandleFunc("/", s.RedirectUI)
	r.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir(fmt.Sprintf("%s", s.conf.StaticDir)))))
	r.NotFoundHandler = http.HandlerFunc(s.NotFoundHandle)

	// 创建HTTP服务
	listenAddr := s.conf.ListenAddr
	if listenAddr == "" {
		listenAddr = ":3005"
	}

	s.server = &http.Server{
		Addr:    listenAddr,
		Handler: r,
	}

	logging.Logger.Infof("Starting HTTP server on http://%s", listenAddr)
	logging.Logger.Infof("Static files served from %s", staticDir)

	// 启动服务
	return s.server.ListenAndServe()
}

// Stop 停止HTTP服务
func (s *Server) Stop() error {
	if s.server == nil {
		return nil
	}

	// 关闭数据库连接
	if s.db != nil {
		_ = s.db.Close()
	}

	logging.Logger.Info("Shutting down HTTP server...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	return s.server.Shutdown(ctx)
}

type AnalyzePDFURLRequest struct {
	URL string `json:"url"`
}

func GetFullURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI)
}

func (s *Server) handleOAuth(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Debugf("Received request to authenticate")
	defer logging.Logger.Debugf("Finished request to authenticate")

	redirect := r.URL.Query().Get("redirect")

	cfg := service.LoadOIDCConfigFromEnv()
	reqURL := GetFullURL(r)
	cfg.RedirectURL = fmt.Sprintf("%s/callback", strings.TrimRight(reqURL, "/"))
	if redirect != "" {
		cfg.RedirectURL = redirect
	}

	//logging.Logger.Infof("OIDC Config: %+v", cfg)

	service, err := service.NewOIDCService(cfg)
	if err != nil {
		logging.Logger.Errorf("Failed to create OIDC service: %v", err)
		s.ResponseError(err, w)
		return
	}
	redirectURL := service.AuthCodeURL(uuid.New().String())

	http.Redirect(w, r, redirectURL, 302)
}

func (s *Server) handleOAuthCallback(w http.ResponseWriter, r *http.Request) {
	logging.Logger.Debugf("Received request to callback")
	defer logging.Logger.Debugf("Finished request to callback")

	code := r.URL.Query().Get("code")
	service, err := service.NewOIDCService(service.LoadOIDCConfigFromEnv())
	if err != nil {
		logging.Logger.Errorf("Failed to create OIDC service: %v", err)
		s.ResponseError(err, w)
		return
	}
	_, err = service.Exchange(r.Context(), code)
	if err != nil {
		logging.Logger.Errorf("Failed to exchange code for token: %v", err)
	}
	apikey, err := s.apiKeyManager.GenerateTemporaryKey(10 * time.Hour)
	if err != nil {
		logging.Logger.Errorf("Failed to generate API key: %v", err)
		s.ResponseError(err, w)
		return
	}
	s.ResponseJSON(apikey, w)
}

func (s *Server) handleGetApiKeyWithCode(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")

	service, err := service.NewOIDCService(service.LoadOIDCConfigFromEnv())
	if err != nil {
		logging.Logger.Errorf("Failed to create OIDC service: %v", err)
		s.ResponseError(err, w)
		return
	}
	_, err = service.Exchange(r.Context(), code)
	if err != nil {
		logging.Logger.Errorf("Failed to exchange code for token: %v", err)
	}
	apikey, err := s.apiKeyManager.GenerateTemporaryKey(10 * time.Hour)
	if err != nil {
		logging.Logger.Errorf("Failed to generate API key: %v", err)
		s.ResponseError(err, w)
		return
	}
	s.ResponseJSON(apikey, w)
}

func (s *Server) ResponseJSON(data interface{}, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(StdAPIResponse{Status: true, Data: data})
}

func (s *Server) ResponseError(err error, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		_ = json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error(), StdAPIResponse: StdAPIResponse{Status: false}})
	} else {
		_ = json.NewEncoder(w).Encode(StdAPIResponse{Status: true})
	}
}

package main

import (
	"net/http"

	"github.com/sirupsen/logrus"

	"openai-forward/config"
	"openai-forward/logging"
	"openai-forward/proxy"
)

func main() {
	// 加载配置
	cfg, err := config.LoadConfig()
	if err != nil {
		logging.Logger.WithError(err).Fatal("Failed to load config")
	}

	// 创建代理
	reverseProxy := proxy.NewOpenAIProxy(cfg)

	// 设置日志级别
	if cfg.LogLevel == "debug" {
		logging.Logger.SetLevel(logrus.DebugLevel)
	} else {
		logging.Logger.SetLevel(logrus.InfoLevel)
	}

	// 设置处理函数
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		logging.Logger.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL.String(),
		}).Debug("Received request")

		reverseProxy.ServeHTTP(w, r)
	})

	// 启动服务器
	logging.Logger.Infof("Starting proxy server on %s", cfg.ListenAddr)
	if err := http.ListenAndServe(cfg.ListenAddr, nil); err != nil {
		logging.Logger.WithError(err).Fatal("Failed to start server")
	}
}

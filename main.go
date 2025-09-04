package main

import (
	"net/http"
	httpService "openai-forward/http"
	"openai-forward/logging"
	"os"
	"os/signal"
	"runtime"
	"syscall"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	conf := &httpService.HTTPConfig{}
	conf.MarginWithENV()

	logging.Logger.Debug("show config detail:")
	logging.Logger.Debug(conf.ToJSON())

	server := httpService.NewServer(conf)

	// 启动服务在goroutine中
	go func() {
		if err := server.Start(); err != nil && err != http.ErrServerClosed {
			logging.Logger.Errorf("Failed to start HTTP server: %v", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logging.Logger.Info("Shutting down server...")
	if err := server.Stop(); err != nil {
		logging.Logger.Errorf("Server forced to shutdown: %v", err)
	}

	logging.Logger.Info("Server exiting")
}

package logging

import (
	"os"

	"github.com/sirupsen/logrus"
)

// Logger 全局日志记录器
var Logger *logrus.Logger

func init() {
	Logger = logrus.New()
	Logger.SetOutput(os.Stdout)
	Logger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	levelFromENV := os.Getenv("PROXY_LOG_LEVEL")
	switch levelFromENV {
	case "debug":
		Logger.SetLevel(logrus.DebugLevel)
	case "warn":
		Logger.SetLevel(logrus.WarnLevel)
	case "error":
		Logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		Logger.SetLevel(logrus.FatalLevel)
	default:
		Logger.SetLevel(logrus.InfoLevel)
	}
}

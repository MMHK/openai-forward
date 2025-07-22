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
}

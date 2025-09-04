package test

import (
	"bytes"
	"encoding/json"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
	"runtime"
)

var testLogger *logrus.Logger

func init() {
	testLogger = logrus.New()
	testLogger.SetOutput(os.Stdout)
	testLogger.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
}

func loadTestEnv() {

	// 获取当前文件路径
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// 加载 .env 文件
	err := godotenv.Load(filepath.Join(dir, "..", ".env"))
	if err != nil {
		testLogger.Errorf("Error loading .env file")
	}
}

func GetLocalPath(file string) string {
	_, filename, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(filename), file)
}

func ToJSON(target interface{}) string {
	str := new(bytes.Buffer)
	encoder := json.NewEncoder(str)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "    ")
	err := encoder.Encode(target)
	if err != nil {
		return err.Error()
	}

	return str.String()
}

func ToMiniyJSON(target interface{}) string {
	str := new(bytes.Buffer)
	encoder := json.NewEncoder(str)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "")
	err := encoder.Encode(target)
	if err != nil {
		return err.Error()
	}

	return str.String()
}

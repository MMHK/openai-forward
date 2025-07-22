package test

import (
	"github.com/joho/godotenv"
	"path/filepath"
	"runtime"
)

func loadTestEnv() {
	// 获取当前文件路径
	_, filename, _, _ := runtime.Caller(0)
	dir := filepath.Dir(filename)

	// 加载 .env 文件
	_ = godotenv.Load(filepath.Join(dir, "..", ".env"))

}

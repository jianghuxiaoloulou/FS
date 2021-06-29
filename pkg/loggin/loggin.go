package loggin

import (
	"os"
	"path/filepath"

	"github.com/wonderivan/logger"
)

func init() {
	// 获取可执行文件的路径
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	// 配置log
	logger.SetLogger(dir + "/log.json")
}

func Debug(f interface{}, v ...interface{}) {
	logger.Debug(f, v...)
}

func Info(f interface{}, v ...interface{}) {
	logger.Info(f, v...)
}

func Warn(f interface{}, v ...interface{}) {
	logger.Warn(f, v...)
}

func Error(f interface{}, v ...interface{}) {
	logger.Error(f, v...)
}

func Fatal(f interface{}, v ...interface{}) {
	logger.Fatal(f, v...)
}

package logger

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"path/filepath"
)

var log *logrus.Logger

func Init(level, logPath string) {
	log = logrus.New()

	// 设置日志级别
	switch level {
	case "debug":
		log.SetLevel(logrus.DebugLevel)
	case "info":
		log.SetLevel(logrus.InfoLevel)
	case "warn":
		log.SetLevel(logrus.WarnLevel)
	case "error":
		log.SetLevel(logrus.ErrorLevel)
	default:
		log.SetLevel(logrus.InfoLevel)
	}

	// 创建日志目录
	if err := os.MkdirAll(logPath, 0755); err != nil {
		log.Fatalf("Failed to create log directory: %v", err)
	}

	// 设置日志输出文件
	logFile, err := os.OpenFile(filepath.Join(logPath, "app.log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(logFile)
	log.SetFormatter(&logrus.JSONFormatter{})
}

func Debug(args ...interface{}) {
	log.Debug(args...)
}

func Info(args ...interface{}) {
	log.Info(args...)
}

func Warn(args ...interface{}) {
	log.Warn(args...)
}

func Error(args ...interface{}) {
	log.Error(args...)
}

func Infof(format string, args ...interface{}) string {
	msg := fmt.Sprintf(format, args...)
	log.Info(msg)
	return msg
}

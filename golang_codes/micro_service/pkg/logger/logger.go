package logger

import (
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

// Init 初始化日志配置
func Init(level string) {
	// 设置日志级别
	logLevel, err := logrus.ParseLevel(strings.ToLower(level))
	if err != nil {
		logLevel = logrus.InfoLevel
	}
	logrus.SetLevel(logLevel)

	// 设置日志格式
	logrus.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "level",
			logrus.FieldKeyMsg:   "message",
			logrus.FieldKeyFunc:  "function",
		},
	})

	// 设置输出到标准输出
	logrus.SetOutput(os.Stdout)

	// 显示文件名和行号
	logrus.SetReportCaller(true)

	logrus.Info("Logger initialized successfully")
}

// GetLogger 获取带有字段的日志实例
func GetLogger(service string) *logrus.Entry {
	return logrus.WithField("service", service)
}

// WithFields 创建带有多个字段的日志实例
func WithFields(fields logrus.Fields) *logrus.Entry {
	return logrus.WithFields(fields)
}

// WithError 创建带有错误信息的日志实例
func WithError(err error) *logrus.Entry {
	return logrus.WithError(err)
}

// WithField 创建带有单个字段的日志实例
func WithField(key string, value interface{}) *logrus.Entry {
	return logrus.WithField(key, value)
}

// Info 记录信息级别日志
func Info(args ...interface{}) {
	logrus.Info(args...)
}

// Infof 记录格式化信息级别日志
func Infof(format string, args ...interface{}) {
	logrus.Infof(format, args...)
}

// Debug 记录调试级别日志
func Debug(args ...interface{}) {
	logrus.Debug(args...)
}

// Debugf 记录格式化调试级别日志
func Debugf(format string, args ...interface{}) {
	logrus.Debugf(format, args...)
}

// Warn 记录警告级别日志
func Warn(args ...interface{}) {
	logrus.Warn(args...)
}

// Warnf 记录格式化警告级别日志
func Warnf(format string, args ...interface{}) {
	logrus.Warnf(format, args...)
}

// Error 记录错误级别日志
func Error(args ...interface{}) {
	logrus.Error(args...)
}

// Errorf 记录格式化错误级别日志
func Errorf(format string, args ...interface{}) {
	logrus.Errorf(format, args...)
}

// Fatal 记录致命错误日志并退出程序
func Fatal(args ...interface{}) {
	logrus.Fatal(args...)
}

// Fatalf 记录格式化致命错误日志并退出程序
func Fatalf(format string, args ...interface{}) {
	logrus.Fatalf(format, args...)
}

// Panic 记录恐慌日志并触发panic
func Panic(args ...interface{}) {
	logrus.Panic(args...)
}

// Panicf 记录格式化恐慌日志并触发panic
func Panicf(format string, args ...interface{}) {
	logrus.Panicf(format, args...)
}

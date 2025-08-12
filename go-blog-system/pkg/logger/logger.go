package logger

import (
	"io"
	"os"
	"path/filepath"

	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Config 日志配置
type Config struct {
	Level  string
	Format string
	Output string
	File   string
}

// InitLogger 初始化日志
func InitLogger() error {
	config := &Config{
		Level:  viper.GetString("logger.level"),
		Format: viper.GetString("logger.format"),
		Output: viper.GetString("logger.output"),
		File:   viper.GetString("logger.file"),
	}

	// 设置日志级别
	level, err := logrus.ParseLevel(config.Level)
	if err != nil {
		level = logrus.InfoLevel
	}
	logrus.SetLevel(level)

	// 设置日志格式
	if config.Format == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	} else {
		logrus.SetFormatter(&logrus.TextFormatter{
			FullTimestamp: true,
		})
	}

	// 设置日志输出
	var output io.Writer = os.Stdout
	if config.Output == "file" && config.File != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.File)
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return err
		}

		// 打开日志文件
		file, err := os.OpenFile(config.File, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			return err
		}
		output = file
	}

	logrus.SetOutput(output)
	return nil
}
package logger

import (
	"fmt"
	"github.com/YeHeng/go-web-api/pkg/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func InitLogger() {

	if err := os.MkdirAll(config.Config.LogConfig.Folder, 0777); err != nil {
		fmt.Println(err.Error())
	}

	encoder := getEncoder()
	level := zapcore.DebugLevel
	_ = level.Set(config.Config.LogConfig.Level)

	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook())),
		level)

	var logger *zap.Logger

	if gin.Mode() == gin.ReleaseMode {
		logger = zap.New(core)
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.Development())
	}

	defer logger.Sync()
	Logger = logger.Sugar()

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewJSONEncoder(encoderConfig)
}

func hook() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   config.Config.LogConfig.Folder + config.Config.LogConfig.Filename,
		MaxSize:    config.Config.LogConfig.MaxSize,
		MaxBackups: config.Config.LogConfig.MaxBackups,
		MaxAge:     config.Config.LogConfig.MaxAge,
		Compress:   config.Config.LogConfig.Compress,
		LocalTime:  config.Config.LogConfig.LocalTime,
	}
}

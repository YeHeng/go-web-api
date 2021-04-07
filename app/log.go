package app

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.SugaredLogger

func init() {

	if err := os.MkdirAll(Config.LogConfig.Folder, 0777); err != nil {
		fmt.Println(err.Error())
	}

	encoder := getEncoder()
	level := zapcore.DebugLevel
	_ = level.Set(Config.LogConfig.Level)

	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook())),
		level)

	var logger *zap.Logger

	if gin.Mode() == gin.ReleaseMode {
		logger = zap.New(core, zap.Fields(zap.String("serviceName", "webhook")))
	} else {
		logger = zap.New(core, zap.AddCaller(), zap.Development(), zap.Fields(zap.String("serviceName", "go-web-tool")))
	}

	defer logger.Sync()
	Logger = logger.Sugar()

}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func hook() *lumberjack.Logger {
	return &lumberjack.Logger{
		Filename:   Config.LogConfig.Folder + Config.LogConfig.Filename,
		MaxSize:    Config.LogConfig.MaxSize,
		MaxBackups: Config.LogConfig.MaxBackups,
		MaxAge:     Config.LogConfig.MaxAge,
		Compress:   Config.LogConfig.Compress,
		LocalTime:  Config.LogConfig.LocalTime,
	}
}

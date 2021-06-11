package logger

import (
	"fmt"
	"github.com/YeHeng/go-web-api/internal/pkg/factory"
	"github.com/YeHeng/go-web-api/pkg/color"
	"github.com/YeHeng/go-web-api/pkg/config"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	factory.Register("log", &logLifecycle{})
}

type logLifecycle struct {
}

func (m *logLifecycle) Destroy() {
}

var log *zap.Logger

func (m *logLifecycle) Init() {

	fmt.Println(color.Green("* [logging init]"))

	cfg := config.Get().Logger

	if err := os.MkdirAll(cfg.Folder, 0777); err != nil {
		fmt.Println(err.Error())
	}

	encoder := getEncoder()
	level := zapcore.DebugLevel
	_ = level.Set(cfg.Level)

	core := zapcore.NewCore(encoder,
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(hook())),
		level)

	if gin.Mode() == gin.ReleaseMode {
		log = zap.New(core)
	} else {
		log = zap.New(core, zap.AddCaller(), zap.Development())
	}

	defer log.Sync()
	// log = logger.Sugar()
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func hook() *lumberjack.Logger {
	cfg := config.Get().Logger
	return &lumberjack.Logger{
		Filename:   cfg.Folder + cfg.Filename,
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.Compress,
		LocalTime:  cfg.LocalTime,
	}
}

func Get() *zap.Logger {
	return log
}

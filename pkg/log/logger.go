package log

import (
	"fmt"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/frangte/gocommon/pkg/xenv"
)

type (
	Config struct {
		AppName string `env:"LOG_APP_NAME" default:"sinbadgo"`
		Level   string `env:"LOG_LEVEL" default:"info"`
		LogDir  string `env:"LOG_DIR" default:"logs"`
		LogMode string `env:"LOG_MODE" default:"production"`
	}

	logger struct {
		z       *zap.Logger
		appName string
	}
)

const (
	DebugLevel = "debug"
	InfoLevel  = "info"
	ErrorLevel = "error"
	WarnLevel  = "warn"
	PanicLevel = "panic"
	FatalLevel = "fatal"

	appNameLabel = "app_name"
)

// Singleton instance of logger for one application
var (
	ins *logger

	mapLevel = map[string]zapcore.Level{
		DebugLevel: zapcore.DebugLevel,
		InfoLevel:  zapcore.InfoLevel,
		ErrorLevel: zapcore.ErrorLevel,
		WarnLevel:  zapcore.WarnLevel,
	}
)

func init() {
	conf, err := xenv.Loads[Config](nil)
	if err != nil {
		panic(err)
	}

	ins = &logger{
		z:       createLogger(*conf),
		appName: conf.AppName,
	}
}

func Info(msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Info(msg, fields...)
}

func Error(msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Error(msg, fields...)
}

func Debug(msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Debug(msg, fields...)
}

func Warn(msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Warn(msg, fields...)
}

func Panic(msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Panic(msg, fields...)
}

func Fatal(msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Fatal(msg, fields...)
}

func Log(level string, msg string, fields ...Field) {
	fields = append(fields, String(appNameLabel, ins.appName))
	ins.z.Log(mapLevel[level], msg, fields...)
}

func Sync() {
	ins.z.Sync()
}

func createLogger(conf Config) *zap.Logger {
	level := zap.NewAtomicLevelAt(mapLevel[conf.Level])
	var zapCore zapcore.Core
	switch conf.LogMode {
	case "development":
		zapCore = developmentZapCore(level)
	case "production":
		zapCore = productionZapCore(conf.LogDir, level, conf.AppName)
	default:
		zapCore = productionZapCore(conf.LogDir, level, conf.AppName)
	}

	return zap.New(zapCore)
}

func productionZapCore(logdir string, level zap.AtomicLevel, appname string) zapcore.Core {
	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filepath.Join(logdir, fmt.Sprintf("%s.log", appname)),
		MaxSize:    10, // megabytes
		MaxBackups: 3,
		MaxAge:     7, // days
	})
	cfg := zap.NewProductionEncoderConfig()
	cfg.TimeKey = "timestamp"
	cfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(cfg)
	return zapcore.NewTee(
		zapcore.NewCore(encoder, file, level),
	)
}

func developmentZapCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)
	cfg := zap.NewDevelopmentEncoderConfig()
	cfg.EncodeLevel = zapcore.CapitalColorLevelEncoder
	encoder := zapcore.NewConsoleEncoder(cfg)
	return zapcore.NewTee(
		zapcore.NewCore(encoder, stdout, level),
	)
}

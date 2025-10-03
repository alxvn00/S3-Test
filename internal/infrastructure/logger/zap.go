package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"log"
	"os"
)

var (
	global       *zap.SugaredLogger
	defaultLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
)

func Init() {
	logger, err := new(defaultLevel, os.Stdout)
	if err != nil {
		log.Printf("error initializing logger: %v", err)
	}
	global = logger
}

func new(lvl zapcore.LevelEnabler, w io.Writer) (*zap.SugaredLogger, error) {
	if lvl == nil {
		lvl = defaultLevel
	}

	cfg := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	enc := zapcore.NewConsoleEncoder(cfg)
	return zap.New(zapcore.NewCore(enc, zapcore.AddSync(w), lvl)).Sugar(), nil
}

func Close() {
	if err := global.Sync(); err != nil {
		log.Printf("error closing logger: %v", err)
	}
}

func GetLogger() *zap.SugaredLogger {
	return global
}

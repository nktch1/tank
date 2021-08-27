package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/nktch1/tank/internal/config"
)

func BuildLogger(conf *config.Config) *zap.Logger {
	zapCfg := zap.NewProductionConfig()

	zapCfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	if conf.Debug {
		zapCfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	if !conf.JSONLogs {
		zapCfg.Encoding = "console"
	}

	logger, err := zapCfg.Build()
	if err != nil {
		panic(err)
	}

	return logger
}

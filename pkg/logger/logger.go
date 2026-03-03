package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Log *zap.Logger

func Init(dev bool) error {
	var config zap.Config

	if dev {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		config = zap.NewProductionConfig()
	}

	logger, err := config.Build()

	if err != nil {
		return err
	}

	Log = logger

	return nil
}

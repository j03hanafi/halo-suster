package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/j03hanafi/halo-suster/common/configs"
)

func setConsoleLogger() (zapcore.Core, []zap.Option) {
	writer := zapcore.AddSync(os.Stdout)

	config := zap.NewDevelopmentEncoderConfig()
	config.EncodeLevel = zapcore.CapitalColorLevelEncoder

	encoder := zapcore.NewConsoleEncoder(config)

	logLevel := zap.NewAtomicLevelAt(zap.InfoLevel)

	if configs.Get().API.DebugMode {
		logLevel = zap.NewAtomicLevelAt(zap.DebugLevel)
	}

	options := append([]zap.Option{}, zap.Development(), zap.AddCaller(), zap.AddStacktrace(zap.DPanicLevel))

	return zapcore.NewCore(encoder, writer, logLevel), options
}

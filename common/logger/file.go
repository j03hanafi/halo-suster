package logger

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/j03hanafi/halo-suster/common/configs"
)

const (
	logMaxSize    = 5
	logMaxAge     = 30
	logMaxBackups = 15
	logSize       = 1024 * 1024

	bufferFlushInterval = 2 * time.Second

	samplerFirst      = 100
	samplerThereafter = 100
)

func setFileLogger() (zapcore.Core, []zap.Option) {
	filename := "logs/" + configs.Get().App.Name + ".log"
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   filename,
		MaxSize:    logMaxSize,
		MaxAge:     logMaxAge,
		MaxBackups: logMaxBackups,
		LocalTime:  true,
	})
	bufferedWriter := &zapcore.BufferedWriteSyncer{
		WS:            writer,
		Size:          logSize,
		FlushInterval: bufferFlushInterval,
	}

	config := zap.NewProductionEncoderConfig()
	config.TimeKey = ""
	config.CallerKey = ""

	encoder := zapcore.NewJSONEncoder(config)

	logLevel := zap.NewAtomicLevelAt(zap.ErrorLevel)
	options := make([]zap.Option, 0)

	core := zapcore.NewCore(encoder, bufferedWriter, logLevel)

	return zapcore.NewSamplerWithOptions(core, time.Second, samplerFirst, samplerThereafter), options
}

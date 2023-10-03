package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var zapLog *zap.Logger

func init() {
	var err error
	config := zap.NewProductionConfig()
	zapcore.TimeEncoderOfLayout("Jan _2 15:04:05.000000000")
	/*
		enccoderConfig := zap.NewProductionEncoderConfig()
		enccoderConfig.StacktraceKey = "" // to hide stacktrace info
		config.EncoderConfig = enccoderConfig
	*/
	zapLog, err = config.Build(zap.AddCallerSkip(1))
	if err != nil {
		panic(err)
	}
}

func Info(message string, fields ...zapcore.Field) {
	zapLog.Info(message, fields...)
}

func Debug(message string, fields ...zapcore.Field) {
	zapLog.Debug(message, fields...)
}

func Error(message string, fields ...zapcore.Field) {
	zapLog.Error(message, fields...)
}

func Fatal(message string, fields ...zapcore.Field) {
	zapLog.Fatal(message, fields...)
}

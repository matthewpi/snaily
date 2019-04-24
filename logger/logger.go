package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"time"
)

var zapLogger *zap.SugaredLogger

// Initialize initializes a new logger.
func Initialize() error {
	config := zap.NewDevelopmentConfig()

	config.EncoderConfig.StacktraceKey = ""
	config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	config.EncoderConfig.EncodeTime = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
		enc.AppendString(t.Format("15:04:05"))
	}
	config.EncoderConfig.EncodeCaller = nil

	logger, err := config.Build()
	if err != nil {
		return err
	}

	defer logger.Sync()
	zapLogger = logger.Sugar()

	return nil
}

// Get returns the logger instance.
func Get() *zap.SugaredLogger {
	return zapLogger
}

// Err converts an error object into a zap field.
func Err(err error) zap.Field {
	return zap.Error(err)
}

// Debug .
func Debug(args ...interface{}) {
	Get().Debug(args...)
}

// Info .
func Info(args ...interface{}) {
	Get().Info(args...)
}

// Warn .
func Warn(args ...interface{}) {
	Get().Warn(args...)
}

// Error .
func Error(args ...interface{}) {
	Get().Error(args...)
}

// Fatal .
func Fatal(args ...interface{}) {
	Get().Fatal(args...)
}

// Panic .
func Panic(args ...interface{}) {
	Get().Panic(args...)
}

// Debugf .
func Debugf(template string, args ...interface{}) {
	Get().Debugf(template, args...)
}

// Infof .
func Infof(template string, args ...interface{}) {
	Get().Infof(template, args...)
}

// Warnf .
func Warnf(template string, args ...interface{}) {
	Get().Warnf(template, args...)
}

// Errorf .
func Errorf(template string, args ...interface{}) {
	Get().Errorf(template, args...)
}

// Fatalf .
func Fatalf(template string, args ...interface{}) {
	Get().Fatalf(template, args...)
}

// Panicf .
func Panicf(template string, args ...interface{}) {
	Get().Panicf(template, args...)
}

// Debugw .
func Debugw(msg string, keysAndValues ...interface{}) {
	Get().Debugw(msg, keysAndValues...)
}

// Infow .
func Infow(msg string, keysAndValues ...interface{}) {
	Get().Infow(msg, keysAndValues...)
}

// Warnw .
func Warnw(msg string, keysAndValues ...interface{}) {
	Get().Warnw(msg, keysAndValues...)
}

// Errorw .
func Errorw(msg string, keysAndValues ...interface{}) {
	Get().Errorw(msg, keysAndValues...)
}

// Fatalw .
func Fatalw(msg string, keysAndValues ...interface{}) {
	Get().Fatalw(msg, keysAndValues...)
}

// Panicw .
func Panicw(msg string, keysAndValues ...interface{}) {
	Get().Panicw(msg, keysAndValues...)
}

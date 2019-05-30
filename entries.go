package log

import (
	"context"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type severity zapcore.Level

var (
	DebugLevel = severity(zap.DebugLevel)
	InfoLevel  = severity(zap.InfoLevel)
	WarnLevel  = severity(zap.WarnLevel)
	ErrorLevel = severity(zap.ErrorLevel)
	FatalLevel = severity(zap.FatalLevel)
)

func Log(ctx context.Context, lvl severity, typeId, message string) {
	DefaultLogger.Check(zapcore.Level(lvl), message).Write(
		zap.String(propertyBuild, defaultBuild),
		zap.String(propertyServiceName, defaultService),
		zap.String(propertyMessageType, typeId),
		zap.String(propertyRequestId, requestId(ctx)),
	)
	DefaultLogger.Sync()
}

func Error(ctx context.Context, err error, message string) {
	Log(ctx, ErrorLevel, "error", message+": "+err.Error())
}

func Errorf(ctx context.Context, err error, format string, v ...interface{}) {
	Error(ctx, err, fmt.Sprintf(format, v...))
}

func Info(ctx context.Context, typeId, message string) {
	Log(ctx, InfoLevel, typeId, message)
}

func Infof(ctx context.Context, typeId, format string, v ...interface{}) {
	Info(ctx, typeId, fmt.Sprintf(format, v...))
}

func Warn(ctx context.Context, typeId, message string) {
	Log(ctx, WarnLevel, typeId, message)
}

func Warnf(ctx context.Context, typeId, format string, v ...interface{}) {
	Warn(ctx, typeId, fmt.Sprintf(format, v...))
}

func Debug(ctx context.Context, typeId, message string) {
	Log(ctx, DebugLevel, typeId, message)
}

func Debugf(ctx context.Context, typeId, format string, v ...interface{}) {
	Debug(ctx, typeId, fmt.Sprintf(format, v...))
}

func Fatal(ctx context.Context, typeId, message string) {
	Log(ctx, FatalLevel, typeId, message)
}

func Fatalf(ctx context.Context, typeId, format string, v ...interface{}) {
	Fatal(ctx, typeId, fmt.Sprintf(format, v...))
}

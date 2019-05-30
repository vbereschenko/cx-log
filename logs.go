package log

import (
	"context"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
)

const (
	propertyServiceName = "service"
	propertyBuild       = "build"
	propertyMessageType = "message-type"
	propertyRequestId   = "request-id"
)

var (
	DefaultLogger = newLogger()

	defaultService, defaultBuild string

	defaultConfig = zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "severity",
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

	encoders = map[string]zapcore.Encoder{
		"json":    zapcore.NewJSONEncoder(defaultConfig),
		"console": zapcore.NewConsoleEncoder(defaultConfig),
	}

	envEncoder = os.Getenv("LOG_ENCODER_NAME")
)

func Config(service, build string) {
	defaultService, defaultBuild = service, build
}

func requestId(ctx context.Context) string {
	var requestId string
	if rid, found := ctx.Value(propertyRequestId).(string); found {
		requestId = rid
	}
	return requestId
}

func newLogger() *zap.Logger {
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	topicDebugging := zapcore.AddSync(ioutil.Discard)
	topicErrors := zapcore.AddSync(ioutil.Discard)

	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	encoder, found := encoders[envEncoder]
	if !found {
		encoder = encoders["json"]
	}

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, topicErrors, highPriority),
		zapcore.NewCore(encoder, consoleErrors, highPriority),
		zapcore.NewCore(encoder, topicDebugging, lowPriority),
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
	)

	return zap.New(core)
}

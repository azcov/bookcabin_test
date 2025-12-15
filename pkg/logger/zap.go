package logger

import (
	"context"
	"fmt"
	"strings"

	"github.com/azcov/bookcabin_test/pkg/consts"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type zapLogger struct {
	z *zap.Logger
}

// NewZapLogger returns a configured Zap logger instance
func NewZapLogger(c LoggerConfig) Logger {
	cfg := zap.NewDevelopmentConfig()
	if c.Environment == "production" {
		cfg = zap.NewProductionConfig()
	}
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	cfg.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// Set log level
	lvl := zap.InfoLevel
	switch strings.ToLower(strings.TrimSpace(c.Level)) {
	case "debug":
		lvl = zap.DebugLevel
	case "warn", "warning":
		lvl = zap.WarnLevel
	case "error":
		lvl = zap.ErrorLevel
	case "fatal":
		lvl = zap.FatalLevel
	}
	cfg.Level = zap.NewAtomicLevelAt(lvl)

	l, err := cfg.Build(zap.AddCallerSkip(2))
	if err != nil {
		panic(err)
	}

	logger := &zapLogger{z: l}

	// Only set default if not set (thread-safety issue here in init but unlikely in this use case)
	if defaultLogger == nil {
		defaultLogger = logger
	}
	return logger
}

func (zl *zapLogger) Info(msg string, keysAndValues ...any) {
	zl.z.Info(msg, toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) Error(msg string, keysAndValues ...any) {
	zl.z.Error(msg, toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) Warn(msg string, keysAndValues ...any) {
	zl.z.Warn(msg, toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) Debug(msg string, keysAndValues ...any) {
	zl.z.Debug(msg, toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) Fatal(msg string, keysAndValues ...any) {
	zl.z.Fatal(msg, toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) InfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	zl.z.Info(zl.addRequestID(ctx, msg), toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) ErrorContext(ctx context.Context, msg string, keysAndValues ...any) {
	zl.z.Error(zl.addRequestID(ctx, msg), toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) WarnContext(ctx context.Context, msg string, keysAndValues ...any) {
	zl.z.Warn(zl.addRequestID(ctx, msg), toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) DebugContext(ctx context.Context, msg string, keysAndValues ...any) {
	zl.z.Debug(zl.addRequestID(ctx, msg), toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) FatalContext(ctx context.Context, msg string, keysAndValues ...any) {
	zl.z.Fatal(zl.addRequestID(ctx, msg), toZapFields(keysAndValues...)...)
}

func (zl *zapLogger) addRequestID(ctx context.Context, msg string) string {
	if ctx == nil {
		zl.z.Error("ctx is nil")
		return msg
	}
	val := ctx.Value(consts.HeaderRequestID)
	if val == nil {
		zl.z.Error("X-Request-ID is nil")
		val := ctx.Value("request_id")
		if val == nil {
			zl.z.Error("request_id is nil")
			return msg
		}
	}
	str, ok := val.(string)
	if !ok || str == "" {
		zl.z.Error("X-Request-ID is not a string or is empty")
		return msg
	}
	return fmt.Sprintf("[R:%v] %s", str, msg)
}

func toZapFields(kvs ...any) []zap.Field {
	if len(kvs) == 0 {
		return nil
	}
	fields := make([]zap.Field, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		if i+1 >= len(kvs) {
			break // odd number of arguments
		}

		key, ok := kvs[i].(string)
		if !ok {
			key = fmt.Sprint(kvs[i])
		}

		// Attempt to use specialized fields for common types
		switch v := kvs[i+1].(type) {
		case string:
			fields = append(fields, zap.String(key, v))
		case int:
			fields = append(fields, zap.Int(key, v))
		case int64:
			fields = append(fields, zap.Int64(key, v))
		case float64:
			fields = append(fields, zap.Float64(key, v))
		case bool:
			fields = append(fields, zap.Bool(key, v))
		case error:
			fields = append(fields, zap.Error(v))
		default:
			fields = append(fields, zap.Any(key, v))
		}
	}
	return fields
}

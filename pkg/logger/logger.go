package logger

import "context"

type Logger interface {
	Debug(msg string, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Error(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
	Fatal(msg string, keysAndValues ...any)

	DebugContext(ctx context.Context, msg string, keysAndValues ...any)
	InfoContext(ctx context.Context, msg string, keysAndValues ...any)
	ErrorContext(ctx context.Context, msg string, keysAndValues ...any)
	WarnContext(ctx context.Context, msg string, keysAndValues ...any)
	FatalContext(ctx context.Context, msg string, keysAndValues ...any)
}

var (
	defaultLogger Logger
)

func init() {
	// Initialize default logger with development config
	defaultLogger = NewZapLogger(LoggerConfig{
		Level:       "info",
		Environment: "development",
	})
}

// SetLogger sets the global default logger
func SetLogger(l Logger) {
	defaultLogger = l
}

func Debug(msg string, keysAndValues ...any) {
	defaultLogger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...any) {
	defaultLogger.Info(msg, keysAndValues...)
}

func Error(msg string, keysAndValues ...any) {
	defaultLogger.Error(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...any) {
	defaultLogger.Warn(msg, keysAndValues...)
}

func Fatal(msg string, keysAndValues ...any) {
	defaultLogger.Fatal(msg, keysAndValues...)
}

func DebugContext(ctx context.Context, msg string, keysAndValues ...any) {
	defaultLogger.DebugContext(ctx, msg, keysAndValues...)
}

func InfoContext(ctx context.Context, msg string, keysAndValues ...any) {
	defaultLogger.InfoContext(ctx, msg, keysAndValues...)
}

func ErrorContext(ctx context.Context, msg string, keysAndValues ...any) {
	defaultLogger.ErrorContext(ctx, msg, keysAndValues...)
}

func WarnContext(ctx context.Context, msg string, keysAndValues ...any) {
	defaultLogger.WarnContext(ctx, msg, keysAndValues...)
}

func FatalContext(ctx context.Context, msg string, keysAndValues ...any) {
	defaultLogger.FatalContext(ctx, msg, keysAndValues...)
}

package log

import (
	"context"
	"os"
)

var DefaultLogger *Logger = New(os.Stderr, WARN)

func WithContextFields(ctx context.Context, v ...any) context.Context {
	fields, ok := ctx.Value(defaultLoggerKey).(*LoggerFields)
	if !ok {
		fields = &LoggerFields{}
		ctx = context.WithValue(ctx, defaultLoggerKey, fields)
	}
	fields.payload = append(fields.payload, v...)
	return ctx
}

func Trace(args ...any) {
	DefaultLogger.output(DefaultLogger.c, TRACE, args...)
}

func TraceContext(ctx context.Context, args ...any) {
	DefaultLogger.output(ctx, TRACE, args...)
}

func Debug(args ...any) {
	DefaultLogger.output(DefaultLogger.c, DEBUG, args...)
}

func DebugContext(ctx context.Context, args ...any) {
	DefaultLogger.output(ctx, DEBUG, args...)
}

func Info(args ...any) {
	DefaultLogger.output(DefaultLogger.c, INFO, args...)
}

func InfoContext(ctx context.Context, args ...any) {
	DefaultLogger.output(ctx, INFO, args...)
}

func Warn(args ...any) {
	DefaultLogger.output(DefaultLogger.c, WARN, args...)
}

func WarnContext(ctx context.Context, args ...any) {
	DefaultLogger.output(ctx, WARN, args...)
}

func Error(args ...any) {
	DefaultLogger.output(DefaultLogger.c, ERROR, args...)
}

func ErrorContext(ctx context.Context, args ...any) {
	DefaultLogger.output(ctx, ERROR, args...)
}

func Fatal(args ...any) {
	DefaultLogger.output(DefaultLogger.c, FATAL, args...)
}

func FatalContext(ctx context.Context, args ...any) {
	DefaultLogger.output(ctx, FATAL, args...)
}

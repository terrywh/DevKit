package infra

import (
	"context"
	"fmt"
	"log/slog"
	"runtime"
)

var DefaultLogger *slog.Logger = slog.Default()

type DefaultLoggerFields struct {
	payload []any
}

type DefaultLoggerKey string

var defaultLoggerKey DefaultLoggerKey = "app.default.log.fields"

func WithContextFields(ctx context.Context, v ...any) context.Context {
	fields, ok := ctx.Value(defaultLoggerKey).(*DefaultLoggerFields)
	if !ok {
		fields = &DefaultLoggerFields{}
		ctx = context.WithValue(ctx, defaultLoggerKey, fields)
	}
	fields.payload = append(fields.payload, v...)
	return ctx
}

func Error(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	DefaultLogger.Error(fmt.Sprint(args...), "file", fmt.Sprintf("%s:%d", file, line))
}

func ErrorContext(ctx context.Context, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	if fields, ok := ctx.Value(defaultLoggerKey).(*DefaultLoggerFields); ok {
		var extra []any
		extra = append(extra, "file", fmt.Sprintf("%s:%d", file, line))
		extra = append(extra, fields.payload...)
		DefaultLogger.ErrorContext(ctx, fmt.Sprint(args...), extra...)
	} else {
		DefaultLogger.ErrorContext(ctx, fmt.Sprint(args...),
			"file", fmt.Sprintf("%s:%d", file, line))
	}
}

func Warn(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	DefaultLogger.Warn(fmt.Sprint(args...), "file", fmt.Sprintf("%s:%d", file, line))
}

func WarnContext(ctx context.Context, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	if fields, ok := ctx.Value(defaultLoggerKey).(*DefaultLoggerFields); ok {
		var extra []any
		extra = append(extra, "file", fmt.Sprintf("%s:%d", file, line))
		extra = append(extra, fields.payload...)
		DefaultLogger.WarnContext(ctx, fmt.Sprint(args...), extra...)
	} else {
		DefaultLogger.WarnContext(ctx, fmt.Sprint(args...),
			"file", fmt.Sprintf("%s:%d", file, line))
	}
}

func Info(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	DefaultLogger.Info(fmt.Sprint(args...), "file", fmt.Sprintf("%s:%d", file, line))
}

func InfoContext(ctx context.Context, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	if fields, ok := ctx.Value(defaultLoggerKey).(*DefaultLoggerFields); ok {
		var extra []any
		extra = append(extra, "file", fmt.Sprintf("%s:%d", file, line))
		extra = append(extra, fields.payload...)
		DefaultLogger.InfoContext(ctx, fmt.Sprint(args...), extra...)
	} else {
		DefaultLogger.InfoContext(ctx, fmt.Sprint(args...),
			"file", fmt.Sprintf("%s:%d", file, line))
	}
}

func Debug(args ...any) {
	_, file, line, _ := runtime.Caller(1)
	DefaultLogger.Debug(fmt.Sprint(args...),
		"file", fmt.Sprintf("%s:%d", file, line))
}

func DebugContext(ctx context.Context, args ...any) {
	_, file, line, _ := runtime.Caller(1)
	if fields, ok := ctx.Value(defaultLoggerKey).(*DefaultLoggerFields); ok {
		var extra []any
		extra = append(extra, "file", fmt.Sprintf("%s:%d", file, line))
		extra = append(extra, fields.payload...)
		DefaultLogger.DebugContext(ctx, fmt.Sprint(args...), extra...)
	} else {
		DefaultLogger.DebugContext(ctx, fmt.Sprint(args...),
			"file", fmt.Sprintf("%s:%d", file, line))
	}
}

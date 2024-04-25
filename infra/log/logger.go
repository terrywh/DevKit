package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"time"
)

type LoggerFields struct {
	payload []any
}

type LoggerKey string

var defaultLoggerKey LoggerKey = "app.default.log.fields"

type Logger struct {
	c context.Context
	w io.Writer
	l int
}

func New(w io.Writer, level int) *Logger {
	return &Logger{
		c: context.Background(),
		w: w,
		l: level,
	}
}

func (l *Logger) timestamp() string {
	return time.Now().Format(time.DateTime)
}

func (l *Logger) source() string {
	_, file, line, _ := runtime.Caller(3)
	return fmt.Sprintf("%s/%s:%d", filepath.Base(filepath.Dir(file)), filepath.Base(file), line)
}

func (l *Logger) GetLevel() int {
	return l.l
}

func (l *Logger) SetLevel(level int) {
	l.l = level
}

func (l *Logger) level(level int) string {
	if level < DEBUG {
		return "TRACE"
	} else if level < INFO {
		return "DEBUG"
	} else if level < WARN {
		return "INFO"
	} else if level < ERROR {
		return "WARN"
	} else {
		return "FATAL"
	}
}

func (l *Logger) output(ctx context.Context, level int, v ...any) {
	if level < l.l {
		return
	}
	var argv []any
	argv = append(argv, l.timestamp())
	argv = append(argv, l.level(level))
	argv = append(argv, l.source())
	if fields, ok := ctx.Value(defaultLoggerKey).(*LoggerFields); ok {
		argv = append(argv, fields.payload...)
	}
	argv = append(argv, "\t")
	argv = append(argv, v...)
	fmt.Fprintln(l.w, argv...)
}

func (l *Logger) Fatal(args ...any) {
	l.output(context.TODO(), FATAL, args...)
	os.Exit(1)
}

func (l *Logger) FatalContext(ctx context.Context, args ...any) {
	l.output(ctx, FATAL, args...)
	os.Exit(1)
}

func (l *Logger) Error(args ...any) {
	l.output(l.c, ERROR, args...)
}

func (l *Logger) ErrorContext(ctx context.Context, args ...any) {
	l.output(ctx, ERROR, args...)
}

func (l *Logger) Warn(args ...any) {
	l.output(l.c, WARN, args...)
}

func (l *Logger) WarnContext(ctx context.Context, args ...any) {
	l.output(l.c, WARN, args...)
}

func (l *Logger) Info(args ...any) {
	l.output(l.c, INFO, args...)
}

func (l *Logger) InfoContext(ctx context.Context, args ...any) {
	l.output(ctx, INFO, args...)
}

func (l *Logger) Debug(args ...any) {
	l.output(l.c, DEBUG, args...)
}

func (l *Logger) DebugContext(ctx context.Context, args ...any) {
	l.output(ctx, DEBUG, args...)
}

func (l *Logger) Trace(args ...any) {
	l.output(l.c, TRACE, args...)
}
func (l *Logger) TraceContext(ctx context.Context, args ...any) {
	l.output(ctx, TRACE, args...)
}

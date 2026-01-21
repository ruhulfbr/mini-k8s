package logger

import (
	"context"
	"log/slog"

	"github.com/labstack/echo/v4"
)

func normalizeContext(v any) context.Context {
	switch c := v.(type) {
	case context.Context:
		return c
	case echo.Context:
		return c.Request().Context()
	default:
		return context.Background()
	}
}

// --------------------
// log helpers
// --------------------

func Info(v any, msg string, args ...any) {
	ctx := normalizeContext(v)
	slog.InfoContext(ctx, msg, "args", args)
}

func Warn(v any, msg string, args ...any) {
	ctx := normalizeContext(v)
	slog.WarnContext(ctx, msg, "args", args)
}

func Debug(v any, msg string, args ...any) {
	ctx := normalizeContext(v)
	slog.DebugContext(ctx, msg, "args", args)
}

func Error(v any, msg string, err error, args ...any) {
	ctx := normalizeContext(v)

	if err != nil {
		args = append(args, slog.Any("error", err))
	}

	slog.ErrorContext(ctx, msg, "args", args)
}

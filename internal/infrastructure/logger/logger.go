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

func Info(v any, msg string, attrs ...any) {
	ctx := normalizeContext(v)
	slog.Log(ctx, slog.LevelInfo, msg, attrs...)
}

func Warn(v any, msg string, attrs ...any) {
	ctx := normalizeContext(v)
	slog.Log(ctx, slog.LevelWarn, msg, attrs...)
}

func Debug(v any, msg string, attrs ...any) {
	ctx := normalizeContext(v)
	slog.Log(ctx, slog.LevelDebug, msg, attrs...)
}

func Error(v any, msg string, err error, attrs ...any) {
	ctx := normalizeContext(v)

	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}

	slog.Log(ctx, slog.LevelError, msg, attrs...)
}

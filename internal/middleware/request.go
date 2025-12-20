package middleware

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

type StructuredLogger struct {
	Logger *slog.Logger
}

func (l *StructuredLogger) NewLogEntry(r *http.Request) middleware.LogEntry {
	var scheme string
	if r.TLS != nil {
		scheme = "https"
	} else {
		scheme = "http"
	}

	// handler := l.Logger.Handler()

	attrs := []slog.Attr{
		slog.String("http_scheme", scheme),
		slog.String("http_proto", r.Proto),
		slog.String("http_method", r.Method),
		slog.String("remote_addr", r.RemoteAddr),
		slog.String("user_agent", r.UserAgent()),
		slog.String("uri", r.RequestURI),
	}

	// Если есть TraceID (например, из OpenTelemetry), добавляем его
	if traceID := r.Header.Get("X-Trace-ID"); traceID != "" {
		attrs = append(attrs, slog.String("trace_id", traceID))
	}

	return &StructuredLoggerEntry{
		Logger: l.Logger,
		Attrs:  attrs,
	}
}

type StructuredLoggerEntry struct {
	Logger *slog.Logger
	Attrs  []slog.Attr
}

func (e *StructuredLoggerEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
    duration := float64(elapsed.Nanoseconds()) / 1_000_000.0 // в миллисекундах

    attrs := append(e.Attrs,
        slog.Int("resp_status", status),
        slog.Int("resp_bytes_length", bytes),
        slog.Float64("resp_elapsed_ms", duration),
    )

    level := slog.LevelInfo
    if status >= 400 {
        level = slog.LevelWarn
    }
    if status >= 500 {
        level = slog.LevelError
    }

    e.Logger.LogAttrs(context.Background(), level, "request completed", attrs...)
}

func (e *StructuredLoggerEntry) Panic(v interface{}, stack []byte) {
	attrs := append(e.Attrs,
		slog.Any("panic", v),
		slog.String("stack", string(stack)),
	)
	e.Logger.LogAttrs(context.Background(), slog.LevelError, "request panic", attrs...)
}

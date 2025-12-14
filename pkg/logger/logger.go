package logger

import (
	"log/slog"
	"os"
	"time"

	"github.com/Maxim-Ba/cv-backend/config"
)

func InitLogger(cfg *config.Config) {
	handler := slog.NewJSONHandler(
		os.Stdout, // или os.Stderr для ошибок
		&slog.HandlerOptions{
			Level:     cfgLvlToSlogLvl(cfg),
			AddSource: true,

			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				if a.Key == slog.TimeKey {
					if t, ok := a.Value.Any().(time.Time); ok {
						a.Value = slog.StringValue(t.Format(time.RFC3339Nano))
					}
				}
				return a
			},
		},
	)
	logger := slog.New(handler)
	logger = logger.With(
		slog.String("env", cfg.AppEnv),
	)
	slog.SetDefault(logger)

	_ = slog.NewLogLogger(handler, slog.LevelInfo)
}

func cfgLvlToSlogLvl(cfg *config.Config) slog.Level {
	switch cfg.LogLevel {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelError
	}
}

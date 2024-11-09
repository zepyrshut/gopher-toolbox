package config

import (
	"fmt"
	"gopher-toolbox/token"
	"io"
	"log/slog"
	"os"
	"time"
)

type App struct {
	DataSource string
	UseCache   bool
	// TODO: Extract Security field and use as ExtendedApp (or think a strategy for better library management)
	Security Security
	AppInfo  AppInfo
}

type AppInfo struct {
	GinMode string
	Version string
}

type Security struct {
	Token     *token.Paseto
	StripeKey string
	Duration  time.Duration
}

func NewLogger(level slog.Level) {
	now := time.Now().Format("2006-01-02")
	if _, err := os.Stat("logs"); os.IsNotExist(err) {
		os.Mkdir("logs", 0755)
	}
	f, _ := os.OpenFile(fmt.Sprintf("logs/log%s.log", now), os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	mw := io.MultiWriter(os.Stdout, f)

	logger := slog.New(slog.NewTextHandler(mw, &slog.HandlerOptions{
		AddSource: true,
		Level:     level,
	}))

	slog.SetDefault(logger)
}

func LogLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelDebug
	}
}

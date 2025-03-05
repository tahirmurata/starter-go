package logger

import (
	"log/slog"
	"os"
	"starter/internal/config"

	"github.com/lmittmann/tint"
)

const (
	LevelFatal slog.Level = 12
)

const (
	ansiReset          = "\u001b[0m"
	ansiBrightRed      = "\u001b[91m"
	ansiBrightRedFaint = "\u001b[91;2m"
	ansiBold           = "\u001b[1m"
)

func init() {
	var logger *slog.Logger

	env := config.New().App.Env

	switch env {
	case "development":
		logger = slog.New(tint.NewHandler(os.Stdout, &tint.Options{
			TimeFormat: "03:04:05",
			Level:      slog.LevelDebug,
			AddSource:  true,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Make errors red
				if err, ok := a.Value.Any().(error); ok {
					aErr := tint.Err(err)
					aErr.Key = a.Key
					return aErr
				}
				// Change ERR+4 to FTL
				if a.Key == slog.LevelKey {
					level := a.Value.Any().(slog.Level)
					if level == LevelFatal {
						a.Value = slog.StringValue(ansiBrightRed + ansiBold + "FTL" + ansiReset)
					}
				}
				return a
			},
		}))
	default:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))
	}

	slog.SetDefault(logger)
}

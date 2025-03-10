package config

import (
	"context"
	"log/slog"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
var (
	k = koanf.New(".")
)

type App struct {
	Port string `koanf:"port"`
	Env  string `koanf:"env"`
}

type Database struct {
	Host     string `koanf:"host"`
	Port     string `koanf:"port"`
	Database string `koanf:"database"`
	Username string `koanf:"username"`
	Password string `koanf:"password"`
	Schema   string `koanf:"schema"`
}

type Config struct {
	App      App      `koanf:"app"`
	Database Database `koanf:"database"`
}

func New() *Config {
	// Load YAML config.
	f := file.Provider("config.toml")
	parser := toml.Parser()
	if err := k.Load(f, parser); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError+4, "Failed to load config", slog.Any("err", err))
		panic(1)
	}

	var configInstance *Config

	// Unmarshal to configInstance
	err := k.Unmarshal("", &configInstance)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError+4, "Failed to unmarshal config", slog.Any("err", err))
		panic(1)
	}

	return configInstance
}

package config

import (
	"context"
	"log/slog"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
var (
	k              = koanf.New(".")
	configInstance *Config
)

type Config struct {
	App struct {
		Port int    `koanf:"port"`
		Env  string `koanf:"env"`
	} `koanf:"app"`
	Database struct {
		Host     string `koanf:"host"`
		Port     int    `koanf:"port"`
		Database string `koanf:"database"`
		Username string `koanf:"username"`
		Password string `koanf:"password"`
		Schema   string `koanf:"schema"`
	} `koanf:"database"`
}

func New() *Config {
	if configInstance != nil {
		return configInstance
	}

	// Load YAML config.
	f := file.Provider("config.yml")
	parser := yaml.Parser()
	if err := k.Load(f, parser); err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError+4, "Failed to load config", slog.Any("err", err))
		panic(1)
	}

	// Unmarshal to configInstance
	err := k.Unmarshal("", &configInstance)
	if err != nil {
		slog.LogAttrs(context.Background(), slog.LevelError+4, "Failed to unmarshal config", slog.Any("err", err))
		panic(1)
	}

	return configInstance
}

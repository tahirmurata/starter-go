package config

import (
	"fmt"
	"sync"

	"github.com/knadh/koanf/parsers/toml/v2"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
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

var once sync.Once
var configInstance *Config

func New() (*Config, error) {
	var configError error
	once.Do(func() {
		k := koanf.New(".")

		// Load YAML config.
		f := file.Provider("config.toml")
		parser := toml.Parser()
		if err := k.Load(f, parser); err != nil {
			configError = fmt.Errorf("loading config: %w", err)
			return
		}

		// Unmarshal to configInstance
		err := k.Unmarshal("", &configInstance)
		if err != nil {
			configError = fmt.Errorf("unmarshaling config: %w", err)
			return
		}
	})

	return configInstance, configError
}

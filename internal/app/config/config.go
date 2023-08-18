package config

import (
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env        string     `yaml:"env" env-default:"dev"`
	HTTPServer HTTPServer `yaml:"httpServer"`
	DB         Database   `yaml:"database"`
	Log        Log        `yaml:"log"`
	Access     Access     `yaml:"access"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"127.0.0.1:9090"`
	Timeout     time.Duration `yaml:"timeout" env-default:"5s"`
	IdleTimeout time.Duration `yaml:"idleTimeout" env-default:"30s"`
}

type Database struct {
	Name             string `yaml:"name" env-default:"test-db"`
	ConnectionString string `yaml:"connectionString" env-default:"mongodb://localhost:27017"`
}

type Log struct {
	Level string `yaml:"level" env-default:"info"`
}

type Access struct {
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true" env:"ACCESS_PASSWORD"`
}

const (
	configPathEnvVar = "CONFIG_PATH"
	passwordEnvVar   = "ACCESS_PASSWORD"
)

func NewConfig(configPath string) (*Config, error) {
	if configPath == "" {
		configPath = os.Getenv(configPathEnvVar)
		if configPath == "" {
			return nil, fmt.Errorf("config path is empty in both argument and env var %s", configPathEnvVar)
		}
	}
	_, err := os.Stat(configPath)
	if err != nil {
		return nil, fmt.Errorf("config file is not found in %s", configPath)
	}
	var cfg Config
	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		return nil, fmt.Errorf("read config file %s: %w", configPath, err)
	}
	return &cfg, nil
}

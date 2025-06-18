package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	Service Service `yaml:"service"`
}

type Service struct {
	Address string `yaml:"address" env:"SERVICE_ADDRESS" env-default:"127.0.0.1:8090"`
	Timeout int    `yaml:"timeout" env:"SERVICE_TIMEOUT" env-default:"10"`
}

func ReadConfig(path string, dotenv ...string) (*Config, error) {
	if err := godotenv.Load(dotenv...); err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
	}

	cfg := new(Config)
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}

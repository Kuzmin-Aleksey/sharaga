package config

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"os"
	"time"
)

type Config struct {
	Log        LogConfig        `yaml:"log"`
	MySQl      MysqlConfig      `yaml:"mysql"`
	HttpServer HttpServerConfig `yaml:"http_server"`
	Redis      RedisConfig      `yaml:"redis"`
	Auth       AuthConfig       `yaml:"auth"`
}

type LogConfig struct {
	File   string `yaml:"file" env:"LOG_FILE" env-default:""`
	Debug  bool   `yaml:"debug" env:"LOG_DEBUG" env-default:"true"`
	Format string `yaml:"format" env:"LOG_FORMAT" env-default:"json"`
}

type MysqlConfig struct {
	Addr           string `yaml:"addr" env:"MYSQL_ADDR" env-default:"localhost:3306"`
	User           string `yaml:"user" env:"MYSQL_USER" env-default:"root"`
	Password       string `yaml:"password" env:"MYSQL_PASSWORD" env-default:""`
	Schema         string `yaml:"schema" env:"MYSQL_SCHEMA" env-default:"public"`
	ConnectTimeout int    `yaml:"connect_timeout" env:"DB_CONNECT_TIMEOUT" env-default:"10"`
}

type RedisConfig struct {
	Host     string `json:"host" yaml:"host"`
	Password string `json:"password" yaml:"password"`
	DB       int    `json:"db" yaml:"db"`
}

type HttpServerConfig struct {
	Addr            string        `yaml:"addr" env:"HTTP_ADDR" env-default:"localhost:8080"`
	ReadTimeout     time.Duration `yaml:"read_timeout" env:"HTTP_READ_TIMEOUT" env-default:"10s"`
	WriteTimeout    time.Duration `yaml:"write_timeout" env:"HTTP_WRITE_TIMEOUT" env-default:"10s"`
	ShutdownTimeout time.Duration `yaml:"shutdown_timeout" env:"HTTP_SHUTDOWN_TIMEOUT" env-default:"10s"`
	Log             HttpLog       `yaml:"log"`
}

type HttpLog struct {
	MaxRequestContentLen   int      `yaml:"max_request_content_len" env:"HTTP_LOG_MAX_REQUEST_CONTENT_LEN" env-default:"2048"`
	MaxResponseContentLen  int      `yaml:"max_response_content_len" env:"HTTP_LOG_MAX_RESPONSE_CONTENT_LEN" env-default:"2048"`
	RequestLoggingContent  []string `yaml:"request_logging_content" env:"HTTP_LOG_REQUEST_LOGGING_CONTENT" env-default:""`
	ResponseLoggingContent []string `yaml:"response_logging_content" env:"HTTP_LOG_RESPONSE_LOGGING_CONTENT" env-default:""`
}

type AuthConfig struct {
	AccessKey       string `yaml:"access_key" env:"ACCESS_KEY" env-default:""`
	AccessTokenTTL  int    `yaml:"access_token_ttl" env:"ACCESS_TOKEN_TTL" env-default:"900"`
	RefreshTokenTTL int    `yaml:"refresh_token_ttl" env:"REFRESH_TOKEN_TTL" env-default:"2592000"`
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

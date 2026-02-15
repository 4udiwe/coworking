package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App      App      `yaml:"app"`
		HTTP     HTTP     `yaml:"http"`
		Postgres Postgres `yaml:"postgres"`
		Log      Log      `yaml:"logger"`
		Auth     Auth     `yaml:"auth"`
		Hasher   Hasher   `yaml:"hasher"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Postgres struct {
		URL            string        `env-required:"true" yaml:"url" env:"POSTGRES_URL"`
		ConnectTimeout time.Duration `env-required:"true" yaml:"connect_timeout" env:"POSTGRES_CONNECT_TIMEOUT"`
	}
	Auth struct {
		AccessTokenSecret  string        `env-required:"true" yaml:"access_token_secret" env:"AUTH_ACCESS_TOKEN_SECRET"`
		RefreshTokenSecret string        `env-required:"true" yaml:"refresh_token_secret" env:"AUTH_REFRESH_TOKEN_SECRET"`
		AccessTokenTTL     time.Duration `env-required:"true" yaml:"access_token_ttl" env:"AUTH_ACCESS_TOKEN_TTL"`
		RefreshTokenTTL    time.Duration `env-required:"true" yaml:"refresh_token_ttl" env:"AUTH_REFRESH_TOKEN_TTL"`
	}
	Hasher struct {
		Cost int `env-required:"true" yaml:"cost" env:"HASHER_COST"`
	}
)

func New(configPath string) (*Config, error) {
	cfg := &Config{}

	if err := cleanenv.ReadConfig(configPath, cfg); err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.ReadConfig: %w", err)
	}

	if err := cleanenv.UpdateEnv(cfg); err != nil {
		return nil, fmt.Errorf("config - NewConfig - cleanenv.UpdateEnv: %w", err)
	}

	return cfg, nil
}

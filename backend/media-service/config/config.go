package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App          App          `yaml:"app"`
		HTTP         HTTP         `yaml:"http"`
		Auth         Auth         `yaml:"auth"`
		Log          Log          `yaml:"logger"`
		MongoDB      MongoDB      `yaml:"mongodb"`
		MinIO        MinIO        `yaml:"minio"`
		Shutdown     Shutdown     `yaml:"shutdown"`
		StaleChecker StaleChecker `yaml:"stale_checker"`
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

	MongoDB struct {
		URI    string `env:"MONGODB_URI" yaml:"uri" env-required:"true"`
		DBName string `env:"MONGODB_DB_NAME" yaml:"name" env-default:"media"`
	}

	MinIO struct {
		Endpoint   string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
		AccessKey  string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
		SecretKey  string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
		BucketName string `yaml:"bucket_name" env:"MINIO_BUCKET_NAME" env-default:"media"`
		UseSSL     bool   `yaml:"use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
	}

	Auth struct {
		PublicKeyPath string `env-required:"true" yaml:"public_key_path" env:"AUTH_PUBLIC_KEY_PATH"`
	}

	Shutdown struct {
		Timeout time.Duration `yaml:"timeout" env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
	}

	StaleChecker struct {
		Interval time.Duration `yaml:"interval" env:"STALE_CHECKER_INTERVAL" env-default:"5m"`
		Limit    int           `yaml:"limit" env:"STALE_CHECKER_LIMIT" env-default:"50"`
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

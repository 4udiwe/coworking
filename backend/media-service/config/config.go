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
		Auth     Auth     `yaml:"auth"`
		Log      Log      `yaml:"logger"`
		GRPC     GRPC     `yaml:"grpc"`
		MongoDB  MongoDB  `yaml:"mongodb"`
		MinIO    MinIO    `yaml:"minio"`
		Shutdown Shutdown `yaml:"shutdown"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	GRPC struct {
		Port string `env-required:"true" yaml:"grpc_port" env:"GRPC_PORT"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	MongoDB struct {
		URI    string `env:"MONGODB_URI" yaml:"uri" env-required:"true"`
		DBName string `env:"MONGODB_DB_NAME" yaml:"name" env-default:"media"`
	}

	MinIO struct {
		Endpoint      string `yaml:"endpoint" env:"MINIO_ENDPOINT" env-default:"localhost:9000"`
		AccessKey     string `yaml:"access_key" env:"MINIO_ACCESS_KEY" env-default:"minioadmin"`
		SecretKey     string `yaml:"secret_key" env:"MINIO_SECRET_KEY" env-default:"minioadmin"`
		BucketName    string `yaml:"bucket_name" env:"MINIO_BUCKET_NAME" env-default:"media"`
		UseSSL        bool   `yaml:"use_ssl" env:"MINIO_USE_SSL" env-default:"false"`
		PublicBaseURL string `yaml:"public_base_url" env:"MINIO_PUBLIC_BASE_URL" env-required:"true"` // пример: "https://media.example.com" или "http://localhost:9000"
	}

	Auth struct {
		PublicKeyPath string `env-required:"true" yaml:"public_key_path" env:"AUTH_PUBLIC_KEY_PATH"`
	}

	Shutdown struct {
		Timeout time.Duration `yaml:"timeout" env:"SHUTDOWN_TIMEOUT" env-default:"30s"`
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

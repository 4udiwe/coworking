package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type (
	Config struct {
		App         App         `yaml:"app"`
		HTTP        HTTP        `yaml:"http"`
		ClickHouse  ClickHouse  `yaml:"clickhouse"`
		Auth        Auth        `yaml:"auth"`
		Log         Log         `yaml:"logger"`
		Kafka       Kafka       `yaml:"kafka"`
		BatchBuffer BatchBuffer `yaml:"batch_buffer"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name" env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	HTTP struct {
		Port string `env-required:"true" yaml:"port" env:"SERVER_PORT"`
	}

	ClickHouse struct {
		Addr string `env-required:"true" yaml:"addr" env:"CLICKHOUSE_ADDR"`
		DB   string `env-required:"true" yaml:"db" env:"CLICKHOUSE_DB"`
		User string `env-required:"true" yaml:"user" env:"CLICKHOUSE_USER"`
		Pass string `env-required:"true" yaml:"pass" env:"CLICKHOUSE_PASS"`
	}

	Log struct {
		Level string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}

	Auth struct {
		PublicKeyPath string `env-required:"true" yaml:"public_key_path" env:"AUTH_PUBLIC_KEY_PATH"`
	}

	Kafka struct {
		Brokers []string `env-required:"true" yaml:"brokers" env:"KAFKA_BROKERS"`
		Topics  struct {
			BookingEvents string `env-required:"true" yaml:"booking_events" env:"KAFKA_BOOKING_EVENTS"`
		} `env-required:"true" yaml:"topics" env:"KAFKA_TOPICS"`
		Consumer KafkaConsumer `yaml:"consumer"`
	}

	KafkaConsumer struct {
		GroupID           string        `yaml:"group_id" env:"KAFKA_CONSUMER_GROUP_ID"`
		MaxWait           time.Duration `yaml:"max_wait" env:"KAFKA_CONSUMER_MAX_WAIT"`
		SessionTimeout    time.Duration `yaml:"session_timeout" env:"KAFKA_CONSUMER_SESSION_TIMEOUT"`
		HeartbeatInterval time.Duration `yaml:"heartbeat_interval" env:"KAFKA_CONSUMER_HEARTBEAT_INTERVAL"`
		CommitInterval    time.Duration `yaml:"commit_interval" env:"KAFKA_CONSUMER_COMMIT_INTERVAL"`
	}

	BatchBuffer struct {
		BatchSize     int           `yaml:"batch_size" env:"BATCH_BUFFER_SIZE"`
		FlushInterval time.Duration `yaml:"flush_interval" env:"BATCH_BUFFER_INTERVAL"`
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

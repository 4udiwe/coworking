package config

import "github.com/ilyakaznacheev/cleanenv"

type Route struct {
	Path     string `yaml:"path"`
	Upstream string `yaml:"upstream"`
}

type Config struct {
	App struct {
		Name string `yaml:"name"`
	} `yaml:"app"`

	HTTP struct {
		Port string `yaml:"port"`
	} `yaml:"http"`

	Logger struct {
		Level string `yaml:"level"`
	} `yaml:"logger"`

	RateLimit struct {
		RequestsPerSecond int `yaml:"requests_per_second"`
	} `yaml:"rate_limit"`

	Routes []Route `yaml:"routes"`
}

func Load(path string) (*Config, error) {
	var cfg Config
	err := cleanenv.ReadConfig(path, &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

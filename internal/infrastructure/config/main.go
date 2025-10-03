package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	S3Client
	HTTPClient
}
type S3Client struct {
	Bucket    string `env:"S3_BUCKET" env-default:"localhost"`
	Endpoint  string `env:"S3_ENDPOINT" env-default:"localhost"`
	Region    string `env:"S3_REGION" env-default:"localhost"`
	AccessKey string `env:"S3_ACCESSKEY" env-default:"localhost"`
	SecretKey string `env:"S3_SECRETKEY" env-default:"localhost"`
}

type HTTPClient struct {
	Host string `env:"APP_HOST" env-default:"localhost"`
	Port string `env:"APP_PORT" env-default:"8080"`
}

func New() (*Config, error) {
	var cfg Config
	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("could not read config: %w", err)
	}
	return &cfg, nil
}

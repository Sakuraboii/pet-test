package config

import "github.com/kelseyhightower/envconfig"

const envPrefix = "QA"

type Config struct {
	DbHost     string `split_words:"true" default:"localhost"`
	DbPort     int    `split_words:"true" default:"5432"`
	DbName     string `split_words:"true" default:"postgres"`
	DbUser     string `split_words:"true" default:"postgres"`
	DbPassword string `split_words:"true" default:"postgres"`
}

func FromEnv() (*Config, error) {
	cfg := &Config{}
	err := envconfig.Process(envPrefix, cfg)
	return cfg, err
}

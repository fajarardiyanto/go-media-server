package model

import (
	"github.com/fajarardiyanto/flt-go-database/interfaces"
	"github.com/fajarardiyanto/flt-go-utils/flags"
)

var cfg = new(Config)

type Config struct {
	Version   string `yaml:"version" default:"v"`
	Name      string `yaml:"name"`
	Port      string `yaml:"port"`
	ApiSecret string `yaml:"api_secret" default:"SECRET"`
	Database  struct {
		Mysql    interfaces.SQLConfig              `yaml:"mysql"`
		Redis    interfaces.RedisProviderConfig    `yaml:"redis"`
		RabbitMQ interfaces.RabbitMQProviderConfig `yaml:"rabbitmq"`
	} `yaml:"database"`
}

func init() {
	flags.Init("config.yaml", cfg)
}

func GetConfig() *Config {
	return cfg
}

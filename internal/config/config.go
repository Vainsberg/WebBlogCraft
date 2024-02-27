package viper

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Сonfigurations struct {
	DbUser       string `yaml:"UserbymySQL" env:"UserbymySQL" env-default:"root"`
	DbPass       string `yaml:"PassbymySQL" env:"PassbymySQL" env-default:""`
	Addr         string `yaml:"Address" env:"Address" env-default:"8080"`
	DbIp         string `yaml:"DbHost" env:"DbHost" env-default:"127.0.0.1"`
	DbPort       string `yaml:"DbPort" env:"DbPort" env-default:"3306"`
	RedisAddr    string `yaml:"RedisAddr" env:"RedisAddr" env-default:"localhost:6379"`
	RedisPass    string `yaml:"RedisPass" env:"RedisPass" env-default:""`
	RedisDb      int    `yaml:"RedisDb" env:"RedisDb" env-default: "0"`
	RABBITMQ_URL string `yaml:"rabbitMQURL" env:"rabbitMQURL" env-default: ""`
}

func NewConfig() (*Сonfigurations, error) {
	var cfg Сonfigurations

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

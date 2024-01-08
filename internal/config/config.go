package viper

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Сonfigurations struct {
	DbUser string `yaml:"UserbymySQL" env:"UserbymySQL" env-default:"root"`
	DbPass string `yaml:"PassbymySQL" env:"PassbymySQL" env-default:""`
	Addr   string `yaml:"Address" env:"Address" env-default:"8080"`
}

func NewConfig() (*Сonfigurations, error) {
	var cfg Сonfigurations

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

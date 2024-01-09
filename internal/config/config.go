package viper

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Сonfigurations struct {
	DbUser string `yaml:"UserbymySQL" env:"UserbymySQL" env-default:"root"`
	DbPass string `yaml:"PassbymySQL" env:"PassbymySQL" env-default:"1111"`
	Addr   string `yaml:"Address" env:"Address" env-default:"8080"`
	DbIp   string `yaml:"IpSQL" env:"IpSQL" env-default:"127.0.0.1"`
	DbPort string `yaml:"PortSQL" env:"PortSQL" env-default:"3306"`
}

func NewConfig() (*Сonfigurations, error) {
	var cfg Сonfigurations

	err := cleanenv.ReadConfig("config.yaml", &cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, nil
}

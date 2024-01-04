package conifg

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	DbUser string
	DbPass string
}

func NewConfig() (*Config, error) {
	var err error

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	dbUser := viper.GetString("UserbymySQL")
	dbPass := viper.GetString("PassbymySQL")

	return &Config{DbUser: dbUser, DbPass: dbPass}, nil
}

package viper

import (
	"fmt"

	"github.com/spf13/viper"
)

type Сonfigurations struct {
	DbUser string
	DbPass string
}

func NewConfig() (*Сonfigurations, error) {
	var err error

	viper.SetConfigFile("config.yaml")
	err = viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	dbUser := viper.GetString("UserbymySQL")
	dbPass := viper.GetString("PassbymySQL")

	return &Сonfigurations{DbUser: dbUser, DbPass: dbPass}, nil
}

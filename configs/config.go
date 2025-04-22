package configs

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func LoadMain() {
	viper.SetConfigFile(".env")
	if err := viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error loading .env file: %s", err.Error())
		panic(err)
	}
	logrus.Debug("Configuration file loaded")
}

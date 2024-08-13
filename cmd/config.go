package main

import (
	"fmt"

	"github.com/spf13/viper"
)

func InitConfig() error {
	viper.AddConfigPath("../configs")
	viper.SetConfigName("config")
	err := viper.ReadInConfig()
	if err != nil {
		return fmt.Errorf("error during reading config")
	}
	return nil
}

package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

func readConfig(configName string) error {
	viper.SetConfigName(configName)
	viper.SetConfigType("yaml")

	path, _ := os.Getwd()
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath(fmt.Sprintf("%s/configs", path))

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}
	return nil
}

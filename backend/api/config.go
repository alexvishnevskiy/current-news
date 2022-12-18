package api

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
)

type Config struct {
	Url        string `mapstructure:"url"`
	Continents struct {
		NorthAmerica []string `mapstructure:"NorthAmerica"`
		Europe       []string `mapstructure:"Europe"`
		Asia         []string `mapstructure:"Asia"`
		Africa       []string `mapstructure:"Africa"`
		SouthAmerica []string `mapstructure:"SouthAmerica"`
	} `mapstructure:"continent"`
	Categories []string `mapstructure:"category"`
}

func (c *Config) GetConf() error {
	viper.SetConfigName("config")
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

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

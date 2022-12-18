package api

import (
	"fmt"
	"github.com/spf13/viper"
	"log"
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

func (c *Config) GetConf() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	path, _ := os.Getwd()
	viper.AddConfigPath(".")
	viper.AddConfigPath("./configs")
	viper.AddConfigPath("../configs")
	viper.AddConfigPath(fmt.Sprintf("%s/configs", path))

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalln("Failed to load config")
	}

	err := viper.Unmarshal(c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

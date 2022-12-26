package config

import (
	"github.com/spf13/viper"
)

type ConfigHead struct {
	Url       string   `mapstructure:"url"`
	ApiKey    string   `mapstructure:"apiKey"`
	Sources   []string `mapstructure:"sources"`
	SetKey    string   `mapstructure:"setKey"`
	Frequency int      `mapstructure:"updateFrequency"`
	PageSize  int      `mapstructure:"pageSize"`
}

func (c *ConfigHead) GetConf() error {
	err := readConfig("config-newsapi")
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

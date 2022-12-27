package config

import "github.com/spf13/viper"

// TODO: add additional checks for data types and etc...

type ConfigArchive struct {
	Url       string `mapstructure:"url"`
	ApiKey    string `mapstructure:"apiKey"`
	SetKey    string `mapstructure:"setKey"`
	StartDate string `mapstructure:"startDate"`
	Frequency string `mapstructure:"frequency"`
}

func (c *ConfigArchive) GetConf() error {
	err := readConfig("config-archive")
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

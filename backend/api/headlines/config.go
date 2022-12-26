package headlines

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

type ConfigAPI struct {
	Url      string   `mapstructure:"url"`
	ApiKey   string   `mapstructure:"apiKey"`
	Sources  []string `mapstructure:"sources"`
	SetKey   string   `mapstructure:"setKey"`
	PageSize int      `mapstructure:"pageSize"`
}

func (c *ConfigAPI) GetConf() error {
	viper.SetConfigName("config-newsapi")
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

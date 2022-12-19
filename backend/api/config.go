package api

import (
	"fmt"
	"github.com/spf13/viper"
	"os"
	"reflect"
)

type Config interface {
	GetConf() error
	GetUrl() string
	GetApiKey() string
	GetContinents() map[string][]string
	GetCategories() []string
}

type ConfigAPI struct {
	Url        string `mapstructure:"url"`
	ApiKey     string `mapstructure:"apiKey"`
	Continents struct {
		NorthAmerica []string `mapstructure:"NorthAmerica"`
		Europe       []string `mapstructure:"Europe"`
		Asia         []string `mapstructure:"Asia"`
		Africa       []string `mapstructure:"Africa"`
		SouthAmerica []string `mapstructure:"SouthAmerica"`
	} `mapstructure:"continent"`
	Categories []string `mapstructure:"category"`
}

type TestConfig struct {
	Url        string
	ApiKey     string
	Continents struct {
		NorthAmerica []string
		Europe       []string
	}
	Categories []string
}

func (c *ConfigAPI) GetConf() error {
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

func (c *ConfigAPI) GetUrl() string {
	return c.Url
}

func (c *ConfigAPI) GetApiKey() string {
	return c.ApiKey
}

func getContinents[K ConfigAPI | TestConfig](c K) map[string][]string {
	continentDict := make(map[string][]string)
	v := reflect.ValueOf(c).Field(2)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		continent := typeOfS.Field(i).Name
		countries := v.Field(i).Interface().([]string)
		continentDict[continent] = countries
	}
	return continentDict
}

func (c *ConfigAPI) GetContinents() map[string][]string {
	return getContinents(*c)
}

func (c *ConfigAPI) GetCategories() []string {
	return c.Categories
}

func (c *TestConfig) GetConf() error {
	var apiConf ConfigAPI
	err := apiConf.GetConf()
	if err != nil {
		return err
	}

	c.CopyConfig(&apiConf)
	return nil
}

func (c *TestConfig) CopyConfig(apiConf *ConfigAPI) {
	c.ApiKey = apiConf.ApiKey
	c.Url = apiConf.Url
	c.Continents = struct {
		NorthAmerica []string
		Europe       []string
	}{NorthAmerica: apiConf.Continents.NorthAmerica, Europe: apiConf.Continents.Europe}
	c.Categories = apiConf.Categories
}

func (c *TestConfig) GetUrl() string {
	return c.Url
}

func (c *TestConfig) GetApiKey() string {
	return c.ApiKey
}

func (c *TestConfig) GetContinents() map[string][]string {
	return getContinents(*c)
}

func (c *TestConfig) GetCategories() []string {
	return c.Categories
}

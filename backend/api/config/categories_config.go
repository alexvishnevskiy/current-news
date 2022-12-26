package config

import (
	"reflect"

	"github.com/spf13/viper"
)

type Config interface {
	GetConf() error
	GetUrl() string
	GetApiKey() string
	GetContinents() map[string][]string
	GetCategories() []string
}

type ConfigCat struct {
	Url        string `mapstructure:"url"`
	ApiKey     string `mapstructure:"apiKey"`
	Continents struct {
		NorthAmerica []string `mapstructure:"NorthAmerica"`
		Europe       []string `mapstructure:"Europe"`
		Asia         []string `mapstructure:"Asia"`
		Africa       []string `mapstructure:"Africa"`
		SouthAmerica []string `mapstructure:"SouthAmerica"`
		Oceania      []string `mapstructure:"Oceania"`
	} `mapstructure:"continent"`
	Categories []string `mapstructure:"category"`
	Frequency  int      `mapstructure:"updateFrequency"`
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

func (c *ConfigCat) GetConf() error {
	err := readConfig("config-newsdata")
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

func (c *ConfigCat) GetUrl() string {
	return c.Url
}

func (c *ConfigCat) GetApiKey() string {
	return c.ApiKey
}

func getContinents[K ConfigCat | TestConfig](c K) map[string][]string {
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

func (c *ConfigCat) GetContinents() map[string][]string {
	return getContinents(*c)
}

func (c *ConfigCat) GetCategories() []string {
	return c.Categories
}

func (c *ConfigCat) GetFrequency() int {
	return c.Frequency
}

func (c *TestConfig) GetConf() error {
	var apiConf ConfigCat
	err := apiConf.GetConf()
	if err != nil {
		return err
	}

	c.CopyConfig(&apiConf)
	return nil
}

func (c *TestConfig) CopyConfig(apiConf *ConfigCat) {
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

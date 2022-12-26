package config

import "github.com/spf13/viper"

type ConfigBackend struct {
	RedisAddr   string `mapstructure:"redis_addr"`
	BackendPort string `mapstructure:"backend_port"`
}

func (c *ConfigBackend) GetConf() error {
	err := readConfig("config-backend")
	if err != nil {
		return err
	}

	err = viper.Unmarshal(c)
	if err != nil {
		return err
	}
	return nil
}

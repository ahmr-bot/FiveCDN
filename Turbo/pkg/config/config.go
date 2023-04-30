package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host                    string
	Port                    int
	ServerName              string
	PoweredBy               string
	WhiteListURL            string
	WhiteListUpdateInterval time.Duration
}

func NewConfig(configFile string) *Config {
	if configFile == "" {
		configFile = "config.toml"
	}
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		Host:                    viper.GetString("server.host"),
		Port:                    viper.GetInt("server.port"),
		ServerName:              viper.GetString("server.name"),
		PoweredBy:               viper.GetString("server.powered_by"),
		WhiteListURL:            viper.GetString("whitelist.url"),
		WhiteListUpdateInterval: viper.GetDuration("whitelist.update_interval"),
	}

}

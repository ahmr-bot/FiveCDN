package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Host                    string
	Port                    int
	ServerName              string
	WhiteListURL            string
	WhiteListUpdateInterval time.Duration
}

func NewConfig(configFile string) *Config {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	return &Config{
		Host:                    viper.GetString("server.host"),
		Port:                    viper.GetInt("server.port"),
		ServerName:              viper.GetString("server.name"),
		WhiteListURL:            viper.GetString("whitelist.url"),
		WhiteListUpdateInterval: viper.GetDuration("whitelist.update_interval"),
	}
}

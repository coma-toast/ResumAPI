package utils

import (
	"os"

	"github.com/spf13/viper"
)

// config is the configuration struct
type Config struct {
	DBPath        string
	PidFilePath   string
	CachePath     string
	LogFilePath   string
	Port          string
	NowPushAPIKey string
}

func GetConf(configPath string) *Config {
	viper.AddConfigPath(configPath)
	viper.SetConfigName("config")
	err := viper.ReadInConfig()

	if err != nil {
		HandleErr("location", "Config", err)
	}

	conf := &Config{}
	err = viper.Unmarshal(conf)

	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(conf.PidFilePath); os.IsNotExist(err) {
		os.Mkdir(conf.PidFilePath, 0777)
	}
	if _, err := os.Stat(conf.CachePath); os.IsNotExist(err) {
		os.Mkdir(conf.CachePath, 0777)
	}
	if _, err := os.Stat(conf.LogFilePath); os.IsNotExist(err) {
		os.Mkdir(conf.LogFilePath, 0777)
	}

	return conf
}

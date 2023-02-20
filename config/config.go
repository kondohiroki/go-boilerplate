package config

import (
	"fmt"
	"log"
	"sync"

	"github.com/spf13/viper"
)

var config *Config
var m sync.Mutex

type Config struct {
	Discord   Discord    `yaml:"discord"`
	Schedules []Schedule `yaml:"schedules"`
}

type Discord struct {
	Token     string `yaml:"token"`
	ChannelID string `yaml:"channelId"`
}

type Schedule struct {
	Job       string `yaml:"job"`
	Cron      string `yaml:"cron"`
	IsEnabled bool   `yaml:"isEnabled"`
}

func GetConfig() *Config {
	return config
}

func SetConfig(configFile string) {
	m.Lock()
	defer m.Unlock()

	viper.SetConfigFile(configFile)
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("Error getting config file, %s", err)
	}

	err = viper.Unmarshal(&config)
	if err != nil {
		fmt.Println("Unable to decode into struct, ", err)
	}
}

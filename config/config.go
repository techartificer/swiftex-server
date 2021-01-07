package config

import (
	"sync"

	"github.com/spf13/viper"
)

var mu sync.Mutex

// LoadConfig load cofiguration form config file
func LoadConfig() error {
	viper.SetConfigName("config") // name of config file
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./")
	viper.AddConfigPath("../")
	viper.AddConfigPath("../../")

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {
		return err
	}
	LoadMongoDB()
	LoadServer()
	return nil
}

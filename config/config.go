package config

import (
	"sync"

	"github.com/spf13/viper"
)

var mu sync.Mutex

func bindEnvs(envs []string) error {
	for _, v := range envs {
		if err := viper.BindEnv(v); err != nil {
			return err
		}
	}
	return nil
}

// LoadConfig load cofiguration form config file
func LoadConfig() error {
	LoadMongoDB()
	LoadServer()
	LoadJWT()
	LoadRedis()
	return nil
}

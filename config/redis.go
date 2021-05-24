package config

import (
	"github.com/spf13/viper"
)

// Redis holds the Redis configuration
type Redis struct {
	Password string
	Host     string
}

var redis Redis

// GetRedis returns the default Redis configuration
func GetRedis() Redis {
	return redis
}

// LoadRedis loads jwt configuration
func LoadRedis() error {
	mu.Lock()
	defer mu.Unlock()
	envs := []string{"REDIS_HOST", "REDIS_PASWORD"}
	bindEnvs(envs)
	redis = Redis{
		Host:     viper.GetString("REDIS_HOST"),
		Password: viper.GetString("REDIS_PASWORD"),
	}
	return nil
}

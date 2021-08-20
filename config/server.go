package config

import (
	"github.com/spf13/viper"
)

// Server holds the server configuration
type Server struct {
	Name       string
	Host       string
	Port       int
	BcryptCost int
	Env        string
}

var server Server

// GetServer returns the default server configuration
func GetServer() Server {
	return server
}

// LoadServer loads server configuration
func LoadServer() {
	mu.Lock()
	defer mu.Unlock()
	envs := []string{"SERVER_HOST", "SERVER_PORT", "SERVER_NAME", "SERVER_BCRYPT_COST", "SERVER_ENV"}
	bindEnvs(envs)
	server = Server{
		Name:       viper.GetString("SERVER_NAME"),
		Host:       viper.GetString("SERVER_HOST"),
		Port:       viper.GetInt("SERVER_PORT"),
		BcryptCost: viper.GetInt("SERVER_BCRYPT_COST"),
		Env:        viper.GetString("SERVER_ENV"),
	}
}

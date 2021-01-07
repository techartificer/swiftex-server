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

	server = Server{
		Name:       viper.GetString("server.name"),
		Host:       viper.GetString("server.host"),
		Port:       viper.GetInt("server.port"),
		BcryptCost: viper.GetInt("server.bcrypt_cost"),
	}
}

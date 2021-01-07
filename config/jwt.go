package config

import (
	"github.com/spf13/viper"
)

// JWT holds the jwt configuration
type JWT struct {
	Secret     string
	TTL        int
	RefreshTTL int
}

var jwt JWT

// GetJWT returns the default JWT configuration
func GetJWT() JWT {
	return jwt
}

// LoadJWT loads jwt configuration
func LoadJWT() {
	mu.Lock()
	defer mu.Unlock()

	jwt = JWT{
		Secret:     viper.GetString("jwt.secret"),
		TTL:        viper.GetInt("jwt.TTL"),
		RefreshTTL: viper.GetInt("jwt.refreshTTL"),
	}
}

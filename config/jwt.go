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
func LoadJWT() error {
	mu.Lock()
	defer mu.Unlock()
	envs := []string{"JWT_SECRET", "JWT_REFRESH_TTL", "JWT_TTL"}
	bindEnvs(envs)
	jwt = JWT{
		Secret:     viper.GetString("JWT_SECRET"),
		TTL:        viper.GetInt("JWT_TTL"),
		RefreshTTL: viper.GetInt("JWT_REFRESH_TTL"),
	}
	return nil
}

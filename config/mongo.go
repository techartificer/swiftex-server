package config

import (
	"github.com/spf13/viper"
)

// Database holds the database configuration
type Database struct {
	Name string
	URL  string
}

var db Database

// DB returns the default database configuration
func DB() Database {
	return db
}

// LoadMongoDB loads database configuration
func LoadMongoDB() {
	mu.Lock()
	defer mu.Unlock()
	envs := []string{"MONGO_DB_NAME", "MONGO_URL"}
	bindEnvs(envs)
	db = Database{
		Name: viper.GetString("MONGO_DB_NAME"),
		URL:  viper.GetString("MONGO_URL"),
	}
}

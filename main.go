package main

import (
	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/database"
)

func init() {
	err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	err = database.ConnectMongo()
	if err != nil {
		panic(err)
	}
}

func main() {
	defer database.DisconnectMongo()
}

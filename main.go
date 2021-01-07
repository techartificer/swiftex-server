package main

import (
	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/logger"
)

func init() {
	logger.SetupLog()
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

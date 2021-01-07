package main

import (
	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/server"
)

func init() {
	logger.SetupLog()

	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	if err := database.ConnectMongo(); err != nil {
		panic(err)
	}
	server.Start()
}

func main() {
	defer database.DisconnectMongo()
}

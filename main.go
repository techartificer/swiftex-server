package main

import (
	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/firebase"
	"github.com/techartificer/swiftex/lib/random"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/models"
	"github.com/techartificer/swiftex/server"
	"github.com/techartificer/swiftex/validators"
)

func init() {
	logger.SetupLog()
	random.Initialize()
	validators.InitValidator()
	if err := config.LoadConfig(); err != nil {
		panic(err)
	}
	if err := database.ConnectMongo(); err != nil {
		panic(err)
	}
	if err := models.InitializeIndex(database.GetDB()); err != nil {
		panic(err)
	}
	if err := firebase.Initialize(); err != nil {
		panic(err)
	}
}

func main() {
	defer func() {
		err := database.DisconnectMongo()
		if err != nil {
			logger.Log.Errorln(err)
		}
	}()
	server.Start()
	//! Don't write code here
}

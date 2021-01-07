package helper

import (
	"time"

	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateAdmin create new admin
func CreateAdmin() {
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	admin := &models.Admin{
		ID:        primitive.NewObjectID(),
		Name:      "Bodda",
		Phone:     "8801710027639",
		Password:  "sadasdas", // TODO: bcrypt
		Email:     "ss@techartificer.com",
		Status:    constants.Active,
		Role:      constants.SuperAdmin,
		CreatedAt: time.Now().UTC(),
	}
	if err := adminRepo.Create(db, admin); err != nil {
		logger.Log.Errorln(err)
		return
	}
	logger.Log.Infoln("Admin created successfully")
}

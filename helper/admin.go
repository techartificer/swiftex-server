package helper

import (
	"time"

	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/password"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// CreateAdmin create new admin
func CreateAdmin() {
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	hash, err := password.HashPassword("@sadat642")
	if err != nil {
		logger.Errorln(err)
		return
	}
	admin := &models.Admin{
		ID:        primitive.NewObjectID(),
		Name:      "Super Admin",
		Phone:     "8801710027639",
		Password:  hash,
		Email:     "super@techartificer.com",
		Status:    constants.Active,
		Role:      constants.SuperAdmin,
		CreatedAt: time.Now().UTC(),
	}
	if err := adminRepo.Create(db, admin); err != nil {
		logger.Errorln(err)
		return
	}
	logger.Infoln("Admin created successfully")
}

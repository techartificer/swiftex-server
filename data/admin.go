package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// AdminRepository ...
type AdminRepository interface {
	Create(db *mongo.Database, admin *models.Admin) error
	FindByID(db *mongo.Database, ID primitive.ObjectID) (*models.Admin, error)
	FindByUsername(db *mongo.Database, username string) (*models.Admin, error)
}

type adminRepositoryImpl struct{}

var adminRepository AdminRepository

// NewAdminRepo returns AdminRepository instance
func NewAdminRepo() AdminRepository {
	if adminRepository == nil {
		adminRepository = &adminRepositoryImpl{}
	}
	return adminRepository
}

func (a *adminRepositoryImpl) Create(db *mongo.Database, admin *models.Admin) error {
	adminCollection := db.Collection(admin.CollectionName())
	if _, err := adminCollection.InsertOne(context.Background(), admin); err != nil {
		return err
	}
	return nil
}

func (a *adminRepositoryImpl) FindByID(db *mongo.Database, ID primitive.ObjectID) (*models.Admin, error) {
	admin := models.Admin{}
	adminCollection := db.Collection(admin.CollectionName())
	if err := adminCollection.FindOne(context.Background(), bson.M{"_id": ID}).Decode(&admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

func (a *adminRepositoryImpl) FindByUsername(db *mongo.Database, username string) (*models.Admin, error) {
	admin := models.Admin{}
	adminCollection := db.Collection(admin.CollectionName())
	if err := adminCollection.FindOne(context.Background(), bson.M{"username": username}).Decode(&admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

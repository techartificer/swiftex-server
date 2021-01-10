package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AdminRepository ...
type AdminRepository interface {
	Create(db *mongo.Database, admin *models.Admin) error
	FindByID(db *mongo.Database, ID primitive.ObjectID) (*models.Admin, error)
	FindByUsername(db *mongo.Database, phone string) (*models.Admin, error)
	UpdateAdminByID(db *mongo.Database, data *validators.ReqAdminUpdate, ID string) (*models.Admin, error)
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

func (a *adminRepositoryImpl) FindByUsername(db *mongo.Database, phone string) (*models.Admin, error) {
	admin := models.Admin{}
	adminCollection := db.Collection(admin.CollectionName())
	if err := adminCollection.FindOne(context.Background(), bson.M{"phone": phone}).Decode(&admin); err != nil {
		return nil, err
	}
	return &admin, nil
}

func (a *adminRepositoryImpl) UpdateAdminByID(db *mongo.Database, data *validators.ReqAdminUpdate, ID string) (*models.Admin, error) {
	admin := &models.Admin{}
	adminCollection := db.Collection(admin.CollectionName())
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", _id}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", data}}
	err = adminCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(admin)
	return admin, err
}

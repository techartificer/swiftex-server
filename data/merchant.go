package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type MerchentRepository interface {
	Create(db *mongo.Database, merchant *models.Merchant) error
	FindByPhone(db *mongo.Database, phone string) (*models.Merchant, error)
}

type merchantRepoImpl struct{}

var merchantRepo MerchentRepository

func NewMerchantRepo() MerchentRepository {
	if merchantRepo == nil {
		merchantRepo = &merchantRepoImpl{}
	}
	return merchantRepo
}

func (m *merchantRepoImpl) Create(db *mongo.Database, merchant *models.Merchant) error {
	merchantCollection := db.Collection(merchant.CollectionName())
	if _, err := merchantCollection.InsertOne(context.Background(), merchant); err != nil {
		return err
	}
	return nil
}

func (m *merchantRepoImpl) FindByPhone(db *mongo.Database, phone string) (*models.Merchant, error) {
	merchant := &models.Merchant{}
	merchantCollection := db.Collection(merchant.CollectionName())
	filter := bson.M{"phone": phone}
	err := merchantCollection.FindOne(context.Background(), filter).Decode(merchant)
	if err != nil {
		return nil, err
	}
	return merchant, nil
}

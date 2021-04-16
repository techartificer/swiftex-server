package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MerchentRepository interface {
	Create(db *mongo.Database, merchant *models.Merchant) error
	FindByPhone(db *mongo.Database, phone string) (*models.Merchant, error)
	Merchants(db *mongo.Database, lastID string) (*[]models.Merchant, error)
	UpdateByPhone(db *mongo.Database, phone string, merchant *models.Merchant) (*models.Merchant, error)
}

type merchantRepoImpl struct{}

var merchantRepo MerchentRepository

func NewMerchantRepo() MerchentRepository {
	if merchantRepo == nil {
		merchantRepo = &merchantRepoImpl{}
	}
	return merchantRepo
}
func (m *merchantRepoImpl) Merchants(db *mongo.Database, lastID string) (*[]models.Merchant, error) {
	merchant := models.Merchant{}
	merchantCollection := db.Collection(merchant.CollectionName())
	query := make(bson.M)
	if lastID != "" {
		_lastID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		query["_id"] = bson.M{"$lt": _lastID}
	}
	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(15)
	cursor, err := merchantCollection.Find(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}
	var merchants []models.Merchant
	if err = cursor.All(context.Background(), &merchants); err != nil {
		return nil, err
	}
	return &merchants, nil
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

func (m *merchantRepoImpl) UpdateByPhone(db *mongo.Database, phone string, merchant *models.Merchant) (*models.Merchant, error) {
	merchantCollection := db.Collection(merchant.CollectionName())
	filter := bson.D{{"phone", phone}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}

	updatedMerchant := &models.Merchant{}
	update := bson.M{"$set": merchant}
	err := merchantCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(updatedMerchant)
	return merchant, err
}

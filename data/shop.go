package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopRepository interface {
	Create(db *mongo.Database, shop *models.Shop) error
	ShopsByOwnerId(db *mongo.Database, owner primitive.ObjectID) (*[]models.Shop, error)
}

type shopRepositoryImpl struct{}

var shopRepository ShopRepository

func NewShopRepo() ShopRepository {
	if shopRepository == nil {
		shopRepository = &shopRepositoryImpl{}
	}
	return shopRepository
}

func (s *shopRepositoryImpl) Create(db *mongo.Database, shop *models.Shop) error {
	shopCollection := db.Collection(shop.CollectionName())
	_, err := shopCollection.InsertOne(context.Background(), shop)
	return err
}

func (a *shopRepositoryImpl) ShopsByOwnerId(db *mongo.Database, owner primitive.ObjectID) (*[]models.Shop, error) {
	shop := &models.Shop{}
	shopCollection := db.Collection(shop.CollectionName())
	query := bson.M{"owner": owner}
	cursor, err := shopCollection.Find(context.Background(), query)
	if err != nil {
		return nil, err
	}
	var shops []models.Shop
	if err = cursor.All(context.Background(), &shops); err != nil {
		return nil, err
	}
	return &shops, nil
}

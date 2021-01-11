package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type ShopRepository interface {
	Create(db *mongo.Database, shop *models.Shop) error
}

type ShopRepositoryImpl struct{}

var shopRepository ShopRepository

func NewShopRepo() ShopRepository {
	if shopRepository == nil {
		shopRepository = &ShopRepositoryImpl{}
	}
	return shopRepository
}

func (s *ShopRepositoryImpl) Create(db *mongo.Database, shop *models.Shop) error {
	shopCollection := db.Collection(shop.CollectionName())
	_, err := shopCollection.InsertOne(context.Background(), shop)
	return err
}

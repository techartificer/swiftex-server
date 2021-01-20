package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ShopRepository interface {
	Create(db *mongo.Database, shop *models.Shop) error
	ShopsByOwnerId(db *mongo.Database, owner primitive.ObjectID) (*[]models.Shop, error)
	ShopByID(db *mongo.Database, ID string) (*models.Shop, error)
	UpdateShopByID(db *mongo.Database, ID string, shop *models.Shop) (*models.Shop, error)
	Shops(db *mongo.Database, lastID string, limit int64) (*[]models.Shop, error)
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

func (a *shopRepositoryImpl) UpdateShopByID(db *mongo.Database, ID string, shop *models.Shop) (*models.Shop, error) {
	shopCollection := db.Collection(shop.CollectionName())
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", _id}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", shop}}
	updatedShop := &models.Shop{}
	err = shopCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(updatedShop)
	return updatedShop, err
}

func (a *shopRepositoryImpl) Shops(db *mongo.Database, lastID string, limit int64) (*[]models.Shop, error) {
	shop := &models.Shop{}
	shopCollection := db.Collection(shop.CollectionName())
	opts := options.Find().SetSort(bson.M{"_id": 1}).SetLimit(limit)

	query := bson.M{}
	if lastID != "" {
		id, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		query = bson.M{"_id": bson.M{"$gt": id}}
	}
	cursor, err := shopCollection.Find(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}

	var shops []models.Shop
	if err = cursor.All(context.Background(), &shops); err != nil {
		return nil, err
	}
	return &shops, nil
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

func (a *shopRepositoryImpl) ShopByID(db *mongo.Database, ID string) (*models.Shop, error) {
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	shop := &models.Shop{}
	shopCollection := db.Collection(shop.CollectionName())
	filter := bson.M{"_id": _id}
	err = shopCollection.FindOne(context.Background(), filter).Decode(shop)
	return shop, err
}

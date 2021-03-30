package data

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type ShopRepository interface {
	Create(db *mongo.Database, shop *models.Shop) (*models.Transaction, error)
	ShopsByOwnerId(db *mongo.Database, owner primitive.ObjectID) (*[]models.Shop, error)
	ShopByID(db *mongo.Database, ID string) (*models.Shop, error)
	UpdateShopByID(db *mongo.Database, ID string, shop *models.Shop) (*models.Shop, error)
	Shops(db *mongo.Database, lastID string, limit int64) (*[]models.Shop, error)
	Search(db *mongo.Database, query primitive.M) (*[]models.Shop, error)
}

type shopRepositoryImpl struct{}

var shopRepository ShopRepository

func NewShopRepo() ShopRepository {
	if shopRepository == nil {
		shopRepository = &shopRepositoryImpl{}
	}
	return shopRepository
}

func (s *shopRepositoryImpl) Create(db *mongo.Database, shop *models.Shop) (*models.Transaction, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())
	callBack := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		trx := models.Transaction{
			ID:        primitive.NewObjectID(),
			ShopID:    shop.ID,
			Owner:     shop.Owner,
			Balance:   0,
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
		}
		trxCollection := db.Collection(trx.CollectionName())
		shopCollection := db.Collection(shop.CollectionName())

		if _, err := shopCollection.InsertOne(sessionCtx, shop); err != nil {
			return nil, err
		}
		if _, err := trxCollection.InsertOne(sessionCtx, trx); err != nil {
			return nil, err
		}
		return &trx, nil
	}
	result, err := session.WithTransaction(context.Background(), callBack, txnOpts)
	if err != nil {
		return nil, err
	}
	transaction := models.Transaction{}
	mapstructure.Decode(result, &transaction)
	return &transaction, err
}

func (a *shopRepositoryImpl) Search(db *mongo.Database, query primitive.M) (*[]models.Shop, error) {
	shop := &models.Shop{}
	shopCollection := db.Collection(shop.CollectionName())
	opts := options.Find().SetLimit(10)

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

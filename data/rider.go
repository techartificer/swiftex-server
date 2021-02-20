package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderRepository interface {
	Create(db *mongo.Database, rider *models.Rider) error
	FindByPhone(db *mongo.Database, phone string) (*models.Rider, error)
	FindByID(db *mongo.Database, ID string) (*models.Rider, error)
}

type riderRepoImpl struct{}

var riderRepo RiderRepository

func NewRiderRepo() RiderRepository {
	if riderRepo != nil {
		return riderRepo
	}
	riderRepo = &riderRepoImpl{}
	return riderRepo
}

func (r *riderRepoImpl) Create(db *mongo.Database, rider *models.Rider) error {
	riderCol := db.Collection(rider.CollectionName())
	_, err := riderCol.InsertOne(context.Background(), rider)
	return err
}

func (r *riderRepoImpl) FindByPhone(db *mongo.Database, phone string) (*models.Rider, error) {
	rider := &models.Rider{}
	riderCol := db.Collection(rider.CollectionName())
	filter := bson.M{"phone": phone}
	if err := riderCol.FindOne(context.Background(), filter).Decode(rider); err != nil {
		return nil, err
	}
	return rider, nil
}

func (r *riderRepoImpl) FindByID(db *mongo.Database, ID string) (*models.Rider, error) {
	rider := &models.Rider{}
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	riderCol := db.Collection(rider.CollectionName())
	filter := bson.M{"_id": _id}
	if err := riderCol.FindOne(context.Background(), filter).Decode(rider); err != nil {
		return nil, err
	}
	return rider, nil
}

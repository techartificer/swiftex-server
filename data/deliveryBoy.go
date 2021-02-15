package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeliveryBoyRepository interface{}

type deliveryBoyImpl struct{}

var deliveryBoyRepo DeliveryBoyRepository

func NewDelivaryBoyRepo() DeliveryBoyRepository {
	if deliveryBoyRepo != nil {
		return deliveryBoyRepo
	}
	return deliveryBoyRepo
}

func CreateDelivery(db *mongo.Database, deliveryBoy *models.DeliveryBoy) error {
	deliveryBoyCol := db.Collection(deliveryBoy.CollectionName())
	_, err := deliveryBoyCol.InsertOne(context.Background(), deliveryBoy)
	return err
}

func FindByPhone(db *mongo.Database, phone string) (*models.DeliveryBoy, error) {
	deliveryBoy := &models.DeliveryBoy{}
	deliveryBoyCol := db.Collection(deliveryBoy.CollectionName())
	filter := bson.M{"phone": phone}
	if err := deliveryBoyCol.FindOne(context.Background(), filter).Decode(deliveryBoy); err != nil {
		return nil, err
	}
	return deliveryBoy, nil
}

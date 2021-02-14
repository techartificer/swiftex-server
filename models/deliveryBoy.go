package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DeliveryBoy struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name            string              `bson:"name" json:"name"`
	Password        string              `bson:"password" json:"-"`
	Phone           string              `bson:"phone" json:"phone"`
	Contact         string              `bson:"contact" json:"contact"`
	Remark          string              `bson:"remark" json:"remark"`
	Salary          int                 `bson:"salary" json:"salary"`
	NID             string              `bson:"NID" json:"NID"`
	Address         string              `bson:"address" json:"address"`
	Hub             string              `bson:"hub" json:"hub"`
	CurrentLocation string              `bson:"currentLocation" json:"currentLocation"`
	CreatedAt       time.Time           `bson:"createdAt,omitempty" json:"createdAt"`
	CreatedBy       *primitive.ObjectID `bson:"createdBy,omitempty" json:"-"`
	UpdatedAt       time.Time           `bson:"updatedAt" json:"updatedAt"`
}

func (D DeliveryBoy) CollectionName() string {
	return "deliveryBoys"
}

func initDeliveryBoyIndex(db *mongo.Database) error {
	deliveryBoy := DeliveryBoy{}
	deliveryBoyCol := db.Collection(deliveryBoy.CollectionName())
	if err := createIndex(deliveryBoyCol, bson.M{"phone": 1}, true); err != nil {
		return err
	}
	if err := createIndex(deliveryBoyCol, bson.M{"hub": 1}, false); err != nil {
		return err
	}
	return nil
}

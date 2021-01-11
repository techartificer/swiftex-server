package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Shop holds shops shop data
type Shop struct {
	ID            primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name          string               `bson:"name,omitempty" json:"name"`
	ShopID        string               `bson:"shopId,omitempty" json:"shopId"`
	Phone         string               `bson:"phone,omitempty" json:"phone"`
	Email         string               `bson:"email,omitempty" json:"email"`
	Address       string               `bson:"address,omitempty" json:"address"`
	PickupAddress string               `bson:"pickupAddress,omitempty" json:"pickupAddress"`
	PickupArea    string               `bson:"pickupArea,omitempty" json:"pickupArea"`
	DeliveryZone  string               `bson:"deliveryZone,omitempty" json:"deliveryZone"`
	Coupon        string               `bson:"coupon,omitempty" json:"coupon"`
	Image         string               `bson:"image,omitempty" json:"image,omitempty"`
	Status        string               `bson:"status,omitempty" json:"status"`
	Owner         primitive.ObjectID   `bson:"owner,omitempty" json:"owner"`
	FBPage        string               `bson:"fbPage,omitempty" json:"fbPage"`
	Moderators    []primitive.ObjectID `bson:"moderators,omitempty" json:"moderators"`
	CreatedAt     time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt      time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (s Shop) CollectionName() string {
	return "shops"
}

func initShopIndex(db *mongo.Database) error {
	shop := Shop{}
	shopCol := db.Collection(shop.CollectionName())
	if err := createIndex(shopCol, bson.M{"shopId": 1}, true); err != nil {
		return err
	}
	if err := createIndex(shopCol, bson.M{"owner": 1}, false); err != nil {
		return err
	}
	if err := createIndex(shopCol, bson.M{"phone": 1}, false); err != nil {
		return err
	}
	return nil
}

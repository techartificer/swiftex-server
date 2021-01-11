package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Merchant holds merchants shop data
type Merchant struct {
	ID        primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Name      string               `bson:"name,omitempty" json:"name"`
	Phone     string               `bson:"phone,omitempty" json:"phone"`
	Email     string               `bson:"email,omitempty" json:"email"`
	Shops     []primitive.ObjectID `bson:"shops,omitempty" json:"shops,omitempty"`
	Password  string               `bson:"password,omitempty" json:"-"`
	Status    string               `bson:"status,omitempty" json:"status"`
	CreatedAt time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt  time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (m Merchant) CollectionName() string {
	return "merchants"
}

func initMerchantIndex(db *mongo.Database) error {
	merchant := Merchant{}
	merchantCol := db.Collection(merchant.CollectionName())
	if err := createIndex(merchantCol, bson.M{"phone": 1}, true); err != nil {
		return err
	}
	if err := createIndex(merchantCol, bson.M{"email": 1}, true); err != nil {
		return err
	}
	return nil
}

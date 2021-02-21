package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Rider struct {
	ID              primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name            string              `bson:"name" json:"name"`
	Password        string              `bson:"password" json:"-"`
	Phone           string              `bson:"phone" json:"phone"`
	Contact         string              `bson:"contact" json:"contact"`
	Remark          string              `bson:"remark" json:"remark"`
	Salary          int32               `bson:"salary" json:"salary"`
	NID             string              `bson:"NID" json:"NID"`
	Address         string              `bson:"address" json:"address"`
	Hub             string              `bson:"hub" json:"hub"`
	CurrentLocation string              `bson:"currentLocation" json:"currentLocation"`
	Status          string              `bson:"status,omitempty" json:"status"`
	CreatedAt       time.Time           `bson:"createdAt,omitempty" json:"createdAt"`
	CreatedBy       *primitive.ObjectID `bson:"createdBy,omitempty" json:"-"`
	UpdatedAt       time.Time           `bson:"updatedAt" json:"updatedAt"`
}

func (D Rider) CollectionName() string {
	return "riders"
}

func initRiderIndex(db *mongo.Database) error {
	rider := Rider{}
	riderCol := db.Collection(rider.CollectionName())
	if err := createIndex(riderCol, bson.M{"phone": 1}, true); err != nil {
		return err
	}
	if err := createIndex(riderCol, bson.M{"hub": 1}, false); err != nil {
		return err
	}
	return nil
}

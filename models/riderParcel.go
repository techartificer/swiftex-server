package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type RiderParcel struct {
	ID         primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	RiderID    primitive.ObjectID `bson:"riderId,omitempty" json:"riderId"`
	OrderID    primitive.ObjectID `bson:"orderId,omitempty" json:"orderId"`
	AssignedBy primitive.ObjectID `bson:"assignedBy,omitempty" json:"assignedBy"`
	Status     string             `bson:"status,omitempty" json:"status"`
	CreatedAt  time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt  time.Time          `bson:"updatedAt" json:"updatedAt"`
}

func (r RiderParcel) CollectionName() string {
	return "riderParcels"
}

func initRiderParcelIndex(db *mongo.Database) error {
	riderParcel := RiderParcel{}
	riderParcelCol := db.Collection(riderParcel.CollectionName())
	if err := createIndex(riderParcelCol, bson.M{"riderId": 1}, false); err != nil {
		return err
	}
	return nil
}

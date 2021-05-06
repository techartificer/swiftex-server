package serializer

import "go.mongodb.org/mongo-driver/bson/primitive"

type AllShops struct {
	ID            primitive.ObjectID `bson:"_id" json:"id"`
	Name          string             `bson:"name,omitempty" json:"name"`
	Address       string             `bson:"address,omitempty" json:"address"`
	PickupAddress string             `bson:"pickupAddress,omitempty" json:"pickupAddress"`
	PickupArea    string             `bson:"pickupArea,omitempty" json:"pickupArea"`
}

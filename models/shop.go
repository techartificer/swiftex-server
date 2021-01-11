package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Shop holds merchants shop data
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
	Moderators    []primitive.ObjectID `bson:"moderators,omitempty" json:"moderators"`
	CreatedAt     time.Time            `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt      time.Time            `bson:"updatedAt,omitempty" json:"updatedAt"`
}

package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// OrderStatus holds order status data
type OrderStatus struct {
	ID              primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	Text            string             `bson:"text,omitempty" json:"text"`
	DeleveryBoyID   primitive.ObjectID `bson:"deleveryBoyId,omitempty" json:"deleveryBoy"`
	ShopModeratorID primitive.ObjectID `bson:"shopModeratorId,omitempty" json:"shopModerator"`
	MerchantID      primitive.ObjectID `bson:"merchantId,omitempty" json:"merchant"`
	AdminID         primitive.ObjectID `bson:"adminId,omitempty" json:"admin"`
	Status          string             `bson:"status,omitempty" json:"status"`
	Time            time.Time          `bson:"time,omitempty" json:"time"`
}

// Order holds order data
type Order struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopID                primitive.ObjectID `bson:"shopId,omitempty" json:"shopId"`
	DeliveryBoy           primitive.ObjectID `bson:"deliveryBoy,omitempty" json:"deliveryBoy"`
	ShopModeratorID       primitive.ObjectID `bson:"shopModeratorId,omitempty" json:"shopModerator"`
	MerchantID            primitive.ObjectID `bson:"merchantId,omitempty" json:"merchant"`
	RecipientName         string             `bson:"recipientName,omitempty" json:"recipientName"`
	RecipientPhone        string             `bson:"recipientPhone,omitempty" json:"recipientPhone"`
	RecipientCity         string             `bson:"recipientCity,omitempty" json:"recipientCity"`
	RecipientThana        string             `bson:"recipientThana,omitempty" json:"recipientThana"`
	RecipientArea         string             `bson:"recipientArea,omitempty" json:"recipientArea"`
	RecipientZip          string             `bson:"recipientZip,omitempty" json:"recipientZip"`
	RecipientAddress      string             `bson:"recipientAddress,omitempty" json:"recipientAddress"`
	PackageCode           string             `bson:"packageCode,omitempty" json:"packageCode"`
	PaymentStatus         string             `bson:"paymentStatus,omitempty" json:"paymentStatus"`
	Price                 float64            `bson:"price,omitempty" json:"price"`
	ParcelType            string             `bson:"parcelType,omitempty" json:"parcelType"`
	RequestedDeliveryTime time.Time          `bson:"requestedDeliveryTime,omitempty" json:"requestedDeliveryTime"`
	PickAddress           string             `bson:"pickAddress,omitempty" json:"pickAddress"`
	PickHub               string             `bson:"pickHub,omitempty" json:"pickHub"`
	Comments              string             `bson:"comments,omitempty" json:"comments"`
	NumberOfItems         int                `bson:"numberOfItems,omitempty" json:"numberOfItems"`
	TrackID               string             `bson:"trackId,omitempty" json:"trackId"`
	DeliveryType          string             `bson:"deliveryType,omitempty" json:"deliveryType"`
	Status                []OrderStatus      `bson:"status,omitempty" json:"status"`
	IsCancelled           bool               `bson:"isCancelled,omitempty" json:"isCancelled"`
	DeliveredAt           time.Time          `bson:"deliverdAt,omitempty" json:"deliverdAt"`
	CreatedAt             time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt              time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (o Order) CollectionName() string {
	return "orders"
}

func initOrderIndex(db *mongo.Database) error {
	order := Order{}
	orderCol := db.Collection(order.CollectionName())
	if err := createIndex(orderCol, bson.M{"shopId": 1}, false); err != nil {
		return err
	}
	if err := createIndex(orderCol, bson.M{"deliveryBoy": 1}, false); err != nil {
		return err
	}
	if err := createIndex(orderCol, bson.M{"trackId": 1}, true); err != nil {
		return err
	}
	if err := createIndex(orderCol, bson.M{"pickHub": 1}, false); err != nil {
		return err
	}
	if err := createIndex(orderCol, bson.M{"recipientArea": 1}, false); err != nil {
		return err
	}
	return nil
}

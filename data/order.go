package data

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// OrderStatus holds order status data
type OrderStatus struct {
	Text          string             `bson:"text,omitempty" json:"text"`
	DeleveryBoyID primitive.ObjectID `bson:"deleveryBoyId,omitempty" json:"deleveryBoy"`
	ModeratorID   primitive.ObjectID `bson:"moderatorId,omitempty" json:"moderator"`
	OwnerID       primitive.ObjectID `bson:"ownerId,omitempty" json:"owner"`
	AdminID       primitive.ObjectID `bson:"adminId,omitempty" json:"adminId"`
	Time          time.Time          `bson:"time,omitempty" json:"time"`
}

// Order holds order data
type Order struct {
	ID                    primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopID                primitive.ObjectID `bson:"shopId,omitempty" json:"shopId"`
	DeliveryBoy           primitive.ObjectID `bson:"DeliveryBoy,omitempty" json:"DeliveryBoy"`
	CreatedBy             primitive.ObjectID `bson:"createdBy,omitempty" json:"createdBy"`
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
	Status                []OrderStatus      `bson:"status,omitempty" json:"status"`
	CreatedAt             time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt              time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (o Order) CollectionName() string {
	return "orders"
}

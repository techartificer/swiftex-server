package serializer

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type shop struct {
	Name    string `bson:"name,omitempty" json:"name"`
	ShopID  string `bson:"shopId,omitempty" json:"shopId"`
	Phone   string `bson:"phone,omitempty" json:"phone"`
	Email   string `bson:"email,omitempty" json:"email"`
	Address string `bson:"address,omitempty" json:"address"`
}

type CashOutRequests struct {
	Shop             shop               `bson:"shop,omitempty" json:"shop"`
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopID           primitive.ObjectID `bson:"shopId,omitempty" json:"shopId"`
	Owner            primitive.ObjectID `bson:"owner,omitempty" json:"owner"`
	Balance          float64            `bson:"balance,omitempty" json:"balance"`
	TrxCode          string             `bson:"trxCode,omitempty" json:"-"`
	TrxCodeExpiresAt int64              `bson:"trxCodeExpiresAt,omitempty" json:"trxCodeExpiresAt"`
	Amount           int64              `bson:"amount,omitempty" json:"amount"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

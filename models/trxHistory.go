package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type TrxType string

const (
	IN  TrxType = "In"
	OUT TrxType = "Out"
)

type TrxHistory struct {
	ID          primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Payment     float64             `bson:"payment,omitempty" json:"payment"`
	TrxID       primitive.ObjectID  `bson:"trxId,omitempty" json:"trxId"`
	ShopID      primitive.ObjectID  `bson:"shopId,omitempty" json:"shopId"`
	OrderID     *primitive.ObjectID `bson:"orderId,omitempty" json:"orderId"`
	PaymentType TrxType             `bson:"paymentType,omitempty" json:"paymentType"`
	Remarks     string              `bson:"remarks,omitempty" json:"remarks"`
	CreatedBy   primitive.ObjectID  `bson:"createdBy,omitempty" json:"createdBy"`
	CreatedAt   time.Time           `bson:"createdAt,omitempty" json:"createdAt"`
}

// CollectionName returns name of the models
func (trx TrxHistory) CollectionName() string {
	return "trxHistories"
}

func initTrxHistoryIndex(db *mongo.Database) error {
	trx := TrxHistory{}
	trxCol := db.Collection(trx.CollectionName())
	if err := createIndex(trxCol, bson.M{"shopId": 1}, false); err != nil {
		return err
	}
	return nil
}

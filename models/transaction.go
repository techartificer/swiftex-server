package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type Transaction struct {
	ID               primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	ShopID           primitive.ObjectID `bson:"shopId,omitempty" json:"shopId"`
	Owner            primitive.ObjectID `bson:"owner,omitempty" json:"owner"`
	Balance          float64            `bson:"balance,omitempty" json:"balance"`
	TrxCode          string             `bson:"trxCode,omitempty" json:"-"`
	TrxCodeExpiresAt int64              `bson:"trxCodeExpiresAt,omitempty" json:"-"`
	Amount           int64              `bson:"amount,omitempty" json:"amount"`
	CreatedAt        time.Time          `bson:"createdAt,omitempty" json:"createdAt"`
	UpdatedAt        time.Time          `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (t Transaction) CollectionName() string {
	return "transactions"
}

func initTransactionIndex(db *mongo.Database) error {
	trx := Transaction{}
	trxCol := db.Collection(trx.CollectionName())
	if err := createIndex(trxCol, bson.M{"shopId": 1}, true); err != nil {
		return err
	}
	return nil
}

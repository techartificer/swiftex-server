package data

import (
	"context"
	"sync"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type TransactionRepository interface {
	TransactionByShopId(db *mongo.Database, shopID string) (*map[string]interface{}, error)
}

type transactionRepoImpl struct{}

var (
	create          sync.Once
	transactionRepo TransactionRepository
)

func NewTransactionRepo() TransactionRepository {
	create.Do(func() {
		transactionRepo = &transactionRepoImpl{}
	})
	return transactionRepo
}

func (t *transactionRepoImpl) TransactionByShopId(db *mongo.Database, shopID string) (*map[string]interface{}, error) {
	_shopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}

	errChan := make(chan error, 3)
	defer close(errChan)
	trxChan := make(chan *models.Transaction)
	defer close(trxChan)
	trxHistoryChan := make(chan *[]models.TrxHistory)
	defer close(trxHistoryChan)

	trx := &models.Transaction{}
	trxCollection := db.Collection(trx.CollectionName())
	query := bson.M{"shopId": _shopID}

	go func() {
		err := trxCollection.FindOne(context.Background(), query).Decode(trx)
		errChan <- err
		trxChan <- trx
	}()

	go func() {
		opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(15)
		cursor, err := trxCollection.Find(context.Background(), query, opts)
		errChan <- err
		var trxHistory *[]models.TrxHistory
		err1 := cursor.All(context.Background(), &trxHistory)
		errChan <- err1
		trxHistoryChan <- trxHistory
	}()

	result := map[string]interface{}{
		"transaction":        <-trxChan,
		"transactionHistory": <-trxHistoryChan,
	}
	if err := <-errChan; err != nil {
		return nil, err
	}
	return &result, nil
}

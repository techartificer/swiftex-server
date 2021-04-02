package data

import (
	"context"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type TransactionRepository interface {
	TransactionByShopId(db *mongo.Database, shopID string) (*map[string]interface{}, error)
	AddTrxHistory(db *mongo.Database, trxHistory *models.TrxHistory) (*map[string]interface{}, error)
}

type transactionRepoImpl struct{}

type trxOrder struct {
	Trx   *models.Transaction
	Order *models.Order
}

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

func (t *transactionRepoImpl) AddTrxHistory(db *mongo.Database, trxHistory *models.TrxHistory) (*map[string]interface{}, error) {
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())
	callBack := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		order := models.Order{}
		orderCollection := db.Collection(order.CollectionName())
		query := bson.M{"_id": trxHistory.OrderID}
		if err := orderCollection.FindOne(sessionCtx, query).Decode(&order); err != nil {
			return nil, err
		}
		if order.DeliveredAt != nil {
			return nil, errors.NewError(string(codes.OrderAlreadyDelevired))
		}
		t := time.Now().UTC() // time
		OrderStatus := models.OrderStatus{
			ID:            primitive.NewObjectID(),
			Text:          "Order succefully delevered at your door",
			Status:        constants.Delivered,
			DeleveryBoyID: &trxHistory.CreatedBy,
			AdminID:       &trxHistory.CreatedBy,
			Time:          time.Now().UTC(),
		}
		orderStatusArray := []models.OrderStatus{OrderStatus}
		push := bson.M{"status": bson.M{"$each": orderStatusArray, "$position": 0}}
		update := bson.M{
			"$set": bson.M{
				"deliveredAt":   &t,
				"currentStatus": constants.Delivered,
			},
			"$push": push,
		}
		if err := orderCollection.FindOneAndUpdate(sessionCtx, query, update, &opt).Decode(&order); err != nil {
			return nil, err
		}
		trxHistory.Payment -= order.Charge
		if order.PaymentStatus == constants.COD {
			onePercent := (order.Price / 100) * 1 // calculating COD charge
			trxHistory.Payment -= onePercent
		}
		trx := &models.Transaction{}
		trxCollection := db.Collection(trx.CollectionName())
		filter := bson.M{"shopId": trxHistory.ShopID}
		update = bson.M{
			"$inc": bson.M{"balance": trxHistory.Payment},
			"$set": bson.M{"updatedAt": t},
		}

		if err := trxCollection.FindOneAndUpdate(sessionCtx, filter, update, &opt).Decode(&trx); err != nil {
			if mongo.ErrNoDocuments == err {
				return nil, errors.NewError(string(codes.TransactionNotFound))
			}
			return nil, err
		}
		trxHistory.TrxID = trx.ID
		trxHistory.CreatedAt = t

		trxHistoryCollection := db.Collection(trxHistory.CollectionName())
		if _, err1 := trxHistoryCollection.InsertOne(sessionCtx, trxHistory); err1 != nil {
			return nil, err
		}
		ret := trxOrder{Trx: trx, Order: &order}

		return ret, nil
	}
	result, err := session.WithTransaction(context.Background(), callBack, txnOpts)
	if err != nil {
		return nil, err
	}
	trx := models.Transaction{}
	order := models.Order{}
	mapstructure.Decode(result.(trxOrder).Trx, &trx)
	mapstructure.Decode(result.(trxOrder).Order, &order)
	ret := map[string]interface{}{
		"transaction": trx,
		"history":     trxHistory,
		"order":       order,
	}
	return &ret, nil
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

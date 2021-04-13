package data

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/password"
	"github.com/techartificer/swiftex/lib/random"
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
	GenerateTrxCode(db *mongo.Database, amount int64, shopID string) (*string, error)
	CashOutRequests(db *mongo.Database, lastID string) (*[]bson.M, error)
	CashOut(db *mongo.Database, _createdBy primitive.ObjectID, trxID, trxCode string) (*models.Transaction, error)
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

func (t *transactionRepoImpl) CashOut(db *mongo.Database, _createdBy primitive.ObjectID, trxID, trxCode string) (*models.Transaction, error) {
	_trxID, err := primitive.ObjectIDFromHex(trxID)
	if err != nil {
		return nil, err
	}

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	callBack := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		query := bson.M{"_id": _trxID}
		trx := &models.Transaction{}
		trxCollection := db.Collection(trx.CollectionName())
		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		if err := trxCollection.FindOne(sessionCtx, query).Decode(&trx); err != nil {
			return nil, err
		}
		cashOutAmount := trx.Amount

		ok := password.CheckPasswordHash(trxCode, trx.TrxCode)
		if !ok || cashOutAmount == 0 {
			return nil, errors.NewError(string(codes.InvalidTrxCode))
		}
		if trx.TrxCodeExpiresAt < time.Now().Unix() {
			return nil, errors.NewError(string(codes.TrxCodeExpired))
		}
		if trx.Balance < float64(cashOutAmount) {
			return nil, errors.NewError(string(codes.InsufficientBalance))
		}

		filter := bson.M{"_id": _trxID}
		update := bson.M{
			"$inc": bson.M{"balance": (-1 * cashOutAmount)},
			"$set": bson.M{
				"amount":    0,
				"updatedAt": time.Now().UTC(),
			},
		}
		if err := trxCollection.FindOneAndUpdate(sessionCtx, filter, update, &opt).Decode(&trx); err != nil {
			return nil, err
		}
		trxHistory := models.TrxHistory{
			ID:          primitive.NewObjectID(),
			PaymentType: models.OUT,
			Payment:     float64(cashOutAmount),
			TrxID:       _trxID,
			OrderID:     nil,
			ShopID:      trx.ShopID,
			CreatedBy:   _createdBy,
			CreatedAt:   time.Now().UTC(),
		}
		trxHistoryCollection := db.Collection(trxHistory.CollectionName())
		if _, err1 := trxHistoryCollection.InsertOne(sessionCtx, trxHistory); err1 != nil {
			return nil, err
		}
		return trx, nil
	}
	result, err := session.WithTransaction(context.Background(), callBack, txnOpts)
	if err != nil {
		return nil, err
	}
	trx := models.Transaction{}
	mapstructure.Decode(result, &trx)
	return &trx, nil
}

func (t *transactionRepoImpl) CashOutRequests(db *mongo.Database, lastID string) (*[]bson.M, error) {
	query := make(bson.M)
	query["amount"] = bson.M{"$gt": 0}
	if lastID != "" {
		_lastID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		query["_id"] = bson.M{"$lt": _lastID}
	}
	matchStage := bson.D{{"$match", query}}
	limitStage := bson.D{{"$limit", 10}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "shops"}, {"localField", "shopId"}, {"foreignField", "_id"}, {"as", "shop"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$shop"}, {"preserveNullAndEmptyArrays", false}}}}
	sortStage := bson.D{{"$sort", bson.D{{"updatedAt", -1}}}}

	trxCollection := db.Collection(models.Transaction{}.CollectionName())
	cursor, err := trxCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, limitStage, lookupStage, unwindStage, sortStage})
	if err != nil {
		return nil, err
	}
	var transactions []bson.M
	if err = cursor.All(context.Background(), &transactions); err != nil {
		return nil, err
	}
	return &transactions, nil
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
		if !order.IsAccepted {
			return nil, errors.NewError(string(codes.OrderNotAcceptedYet))
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
	trxChan := make(chan *models.Transaction)
	defer close(trxChan)
	trxHistoryChan := make(chan *[]models.TrxHistory)
	defer close(trxHistoryChan)

	trx := &models.Transaction{}
	trxCollection := db.Collection(trx.CollectionName())
	trxHistoryCollection := db.Collection(models.TrxHistory{}.CollectionName())
	query := bson.M{"shopId": _shopID}

	go func() {
		err := trxCollection.FindOne(context.Background(), query).Decode(trx)
		errChan <- err
		trxChan <- trx
	}()

	go func() {
		opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(15)
		cursor, err := trxHistoryCollection.Find(context.Background(), query, opts)
		errChan <- err
		var result []models.TrxHistory
		err1 := cursor.All(context.Background(), &result)
		errChan <- err1
		trxHistoryChan <- &result
	}()

	result := map[string]interface{}{
		"transaction":        <-trxChan,
		"transactionHistory": <-trxHistoryChan,
	}
	close(errChan)
	for cerr := range errChan {
		if cerr != nil {
			return nil, cerr
		}
	}
	return &result, nil
}

func (t *transactionRepoImpl) GenerateTrxCode(db *mongo.Database, amount int64, shopID string) (*string, error) {
	_shopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}

	trx := &models.Transaction{}
	trxCollection := db.Collection(trx.CollectionName())

	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())
	query := bson.M{"shopId": _shopID}

	callBack := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		if err1 := trxCollection.FindOne(sessionCtx, query).Decode(trx); err1 != nil {
			return nil, err1
		}
		if trx.Balance < float64(amount) {
			return nil, errors.NewError(string(codes.InsufficientBalance))
		}
		expiresAt := time.Now().Local().Add(time.Hour * time.Duration(24*3)).Unix()
		trxCode, err1 := random.GenerateRandomCode(6)
		if err1 != nil {
			log.Println(err1)
			return nil, err1
		}
		hashTrxCode, err := password.HashPassword(trxCode)
		if err != nil {
			return nil, err
		}
		update := bson.M{"$set": bson.M{
			"trxCode":          hashTrxCode,
			"trxCodeExpiresAt": expiresAt,
			"amount":           amount,
			"updatedAt":        time.Now().UTC(),
		}}
		if err1 := trxCollection.FindOneAndUpdate(sessionCtx, query, update, &opt).Decode(trx); err != nil {
			return nil, err1
		}
		return trxCode, nil
	}
	result, err := session.WithTransaction(context.Background(), callBack, txnOpts)
	if err != nil {
		return nil, err
	}
	code := result.(string)
	return &code, err
}

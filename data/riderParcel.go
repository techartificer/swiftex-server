package data

import (
	"context"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type RiderParcelRepository interface {
	Create(db *mongo.Database, parcel *models.RiderParcel) (*models.Order, error)
	ParcelsByRiderId(db *mongo.Database, riderID, lastID string) (*[]bson.M, error)
}

type riderParcelImpl struct{}

var riderParcelRepo RiderParcelRepository

func NewRiderParcelRepo() RiderParcelRepository {
	if riderParcelRepo == nil {
		riderParcelRepo = &riderParcelImpl{}
	}
	return riderParcelRepo
}

func (p riderParcelImpl) Create(db *mongo.Database, parcel *models.RiderParcel) (*models.Order, error) {
	riderParcelCollection := db.Collection(parcel.CollectionName())
	wc := writeconcern.New(writeconcern.WMajority())
	rc := readconcern.Snapshot()
	txnOpts := options.Transaction().SetWriteConcern(wc).SetReadConcern(rc)
	session, err := db.Client().StartSession()
	if err != nil {
		return nil, err
	}
	defer session.EndSession(context.Background())

	callBack := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		if _, err := riderParcelCollection.InsertOne(sessionCtx, parcel); err != nil {
			return nil, err
		}
		status := models.OrderStatus{
			ID:              primitive.NewObjectID(),
			DeleveryBoyID:   &parcel.RiderID,
			ShopModeratorID: nil,
			MerchantID:      nil,
			AdminID:         &parcel.AssignedBy,
			Status:          constants.InTransit,
			Text:            "Rider have picked your parcel",
			Time:            time.Now().UTC(),
		}
		orderStatusArray := []models.OrderStatus{status}
		push := bson.M{"status": bson.M{"$each": orderStatusArray, "$position": 0}}
		filter := bson.M{"_id": parcel.OrderID}
		after := options.After
		opt := options.FindOneAndUpdateOptions{
			ReturnDocument: &after,
		}
		updatedOrder := models.Order{}
		orderCollection := db.Collection(updatedOrder.CollectionName())
		query := bson.M{"$set": bson.M{"currentStatus": constants.InTransit}, "$push": push}
		err = orderCollection.FindOneAndUpdate(context.Background(), filter, query, &opt).Decode(&updatedOrder)
		if err != nil {
			return nil, err
		}
		return &updatedOrder, nil
	}

	result, err := session.WithTransaction(context.Background(), callBack, txnOpts)
	if err != nil {
		return nil, err
	}

	order := models.Order{}
	mapstructure.Decode(result, &order)
	return &order, nil
}

func (p riderParcelImpl) ParcelsByRiderId(db *mongo.Database, riderID, lastID string) (*[]bson.M, error) {
	query := make(bson.M)
	_riderID, err := primitive.ObjectIDFromHex(riderID)
	if err != nil {
		return nil, err
	}
	query["riderId"] = _riderID

	if lastID != "" {
		_lastID, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			return nil, err
		}
		query["_id"] = bson.M{"$lt": _lastID}
	}

	riderParcelCollection := db.Collection(models.RiderParcel{}.CollectionName())

	matchStage := bson.D{{"$match", query}}
	lookupStage := bson.D{{"$lookup", bson.D{{"from", "orders"}, {"localField", "orderId"}, {"foreignField", "_id"}, {"as", "order"}}}}
	unwindStage := bson.D{{"$unwind", bson.D{{"path", "$order"}, {"preserveNullAndEmptyArrays", false}}}}
	sortStage := bson.D{{"$sort", bson.D{{"_id", -1}}}}

	cursor, err := riderParcelCollection.Aggregate(context.Background(), mongo.Pipeline{matchStage, lookupStage, unwindStage, sortStage})
	if err != nil {
		return nil, err
	}

	var parcels []bson.M
	if err = cursor.All(context.Background(), &parcels); err != nil {
		return nil, err
	}
	return &parcels, nil
}

package data

import (
	"context"
	"time"

	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/lib/helper"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type OrderRepository interface {
	Create(db *mongo.Database, order *models.Order) error
	Orders(db *mongo.Database, query primitive.M) (*[]models.Order, error)
	UpdateOrder(db *mongo.Database, order *models.Order, ID, shopID string) (*models.Order, error)
	AddOrderStatus(db *mongo.Database, orderStatus *models.OrderStatus, ID string) (*models.Order, error)
	OrderByID(db *mongo.Database, ID string) (*models.Order, error)
	TrackOrder(db *mongo.Database, trackID string) (*models.Order, error)
	Dashboard(db *mongo.Database, shopID string, startDate, endDate *time.Time) (map[string]int64, error)
}

type orderRepositoryImpl struct{}

var orderRepository OrderRepository

func NewOrderRepo() OrderRepository {
	if orderRepository == nil {
		orderRepository = &orderRepositoryImpl{}
	}
	return orderRepository
}

func (o *orderRepositoryImpl) Dashboard(db *mongo.Database, shopID string, startDate, endDate *time.Time) (map[string]int64, error) {
	order := &models.Order{}
	orderCollection := db.Collection(order.CollectionName())
	_shopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}
	query := make(bson.M)
	query["shopId"] = _shopID
	if !startDate.IsZero() && !endDate.IsZero() {
		query["$and"] = []bson.M{
			{"createdAt": bson.M{"$gte": startDate}},
			{"createdAt": bson.M{"$lte": endDate}},
		}
	}
	totalChan := make(chan int64)
	defer close(totalChan)
	errChan := make(chan error, 5)
	defer close(errChan)
	go func() {
		cnt, err1 := orderCollection.CountDocuments(context.Background(), query)
		errChan <- err1
		totalChan <- cnt
	}()
	deliveredChan := make(chan int64)
	defer close(deliveredChan)
	go func() {
		query2 := helper.CopyMap(query)
		query2["deliverdAt"] = bson.M{"$gt": time.Time{}}
		cnt, err1 := orderCollection.CountDocuments(context.Background(), query2)
		errChan <- err1
		deliveredChan <- cnt
	}()
	transitChan := make(chan int64)
	defer close(transitChan)
	go func() {
		query2 := helper.CopyMap(query)
		query2["currentStatus"] = constants.InTransit
		cnt, err1 := orderCollection.CountDocuments(context.Background(), query2)
		errChan <- err1
		transitChan <- cnt
	}()
	returnedChan := make(chan int64)
	defer close(returnedChan)
	go func() {
		query2 := helper.CopyMap(query)
		query2["currentStatus"] = constants.Returned
		cnt, err1 := orderCollection.CountDocuments(context.Background(), query2)
		errChan <- err1
		returnedChan <- cnt
	}()
	data := make(map[string]int64)
	data["total"] = <-totalChan
	data["delivered"] = <-deliveredChan
	data["returned"] = <-returnedChan
	data["inTransit"] = <-transitChan

	err = <-errChan
	if err != nil {
		return nil, err
	}
	return data, err
}

func (o *orderRepositoryImpl) Create(db *mongo.Database, order *models.Order) error {
	orderCollection := db.Collection(order.CollectionName())
	_, err := orderCollection.InsertOne(context.Background(), order)
	return err
}

func (o *orderRepositoryImpl) TrackOrder(db *mongo.Database, trackID string) (*models.Order, error) {
	logger.Log.Println("trackId: ", trackID)
	order := &models.Order{}
	orderCollection := db.Collection(order.CollectionName())
	filter := bson.M{"trackId": trackID}

	err := orderCollection.FindOne(context.Background(), filter).Decode(order)
	return order, err
}

func (o *orderRepositoryImpl) OrderByID(db *mongo.Database, ID string) (*models.Order, error) {
	order := &models.Order{}
	orderCollection := db.Collection(order.CollectionName())
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", _id}}

	err = orderCollection.FindOne(context.Background(), filter).Decode(order)
	return order, err
}

func (o *orderRepositoryImpl) Orders(db *mongo.Database, query primitive.M) (*[]models.Order, error) {
	order := models.Order{}
	orderCollection := db.Collection(order.CollectionName())

	opts := options.Find().SetSort(bson.M{"_id": -1}).SetLimit(15)
	cursor, err := orderCollection.Find(context.Background(), query, opts)
	if err != nil {
		return nil, err
	}

	var orders []models.Order
	if err = cursor.All(context.Background(), &orders); err != nil {
		return nil, err
	}
	return &orders, nil
}

func (o *orderRepositoryImpl) UpdateOrder(db *mongo.Database, order *models.Order, ID, shopID string) (*models.Order, error) {
	orderCollection := db.Collection(order.CollectionName())
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	_shopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", _id}, {"shopId", _shopID}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", order}}
	updatedOrder := &models.Order{}
	err = orderCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(&updatedOrder)
	return updatedOrder, err
}

func (o *orderRepositoryImpl) AddOrderStatus(db *mongo.Database, orderStatus *models.OrderStatus, ID string) (*models.Order, error) {
	updatedOrder := &models.Order{}
	orderCollection := db.Collection(updatedOrder.CollectionName())
	_id, err := primitive.ObjectIDFromHex(ID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{{"_id", _id}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	query := make(bson.M)
	if orderStatus.Status == constants.Accepted {
		query["isAccepted"] = true
	}
	if orderStatus.Status == constants.Declined {
		query["isCancelled"] = true
	}
	if orderStatus.Status == constants.Delivered {
		query["deliveredAt"] = time.Now().UTC()
	}
	query["currentStatus"] = orderStatus.Status
	orderStatusArray := []models.OrderStatus{*orderStatus}
	push := bson.M{"status": bson.M{"$each": orderStatusArray, "$position": 0}}
	err = orderCollection.FindOneAndUpdate(context.Background(), filter, bson.M{"$set": query, "$push": push}, &opt).Decode(updatedOrder)
	return updatedOrder, err
}

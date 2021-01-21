package data

import (
	"context"

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
}

type orderRepositoryImpl struct{}

var orderRepository OrderRepository

func NewOrderRepo() OrderRepository {
	if orderRepository == nil {
		orderRepository = &orderRepositoryImpl{}
	}
	return orderRepository
}

func (o *orderRepositoryImpl) Create(db *mongo.Database, order *models.Order) error {
	orderCollection := db.Collection(order.CollectionName())
	_, err := orderCollection.InsertOne(context.Background(), order)
	return err
}

func (o *orderRepositoryImpl) Orders(db *mongo.Database, query primitive.M) (*[]models.Order, error) {
	order := models.Order{}
	orderCollection := db.Collection(order.CollectionName())

	opts := options.Find().SetSort(bson.M{"_id": 1}).SetLimit(15)
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
	err = orderCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(updatedOrder)
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
	update := bson.M{"$push": bson.M{"status": orderStatus}}
	err = orderCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(updatedOrder)
	return updatedOrder, err
}

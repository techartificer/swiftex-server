package data

import (
	"context"

	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/mongo"
)

type OrderRepository interface {
	Create(db *mongo.Database, order *models.Order) error
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

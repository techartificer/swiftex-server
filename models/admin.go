package models

import (
	"time"

	"github.com/techartificer/swiftex/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// Admin model holds the admin's data
type Admin struct {
	ID         primitive.ObjectID  `bson:"_id,omitempty" json:"id"`
	Name       string              `bson:"name,omitempty" json:"name"`
	Phone      string              `bson:"phone,omitempty" json:"phone"`
	Email      string              `bson:"email,omitempty" json:"email"`
	Password   string              `bson:"password,omitempty" json:"-"`
	ProfilePic string              `bson:"profilePic,omitempty" json:"profilePic,omitempty"`
	Status     string              `bson:"status,omitempty" json:"status"`
	Role       constants.AdminRole `bson:"role,omitempty" json:"role"`
	CreatedAt  time.Time           `bson:"createdAt,omitempty" json:"createdAt"`
	UpdateAt   time.Time           `bson:"updatedAt,omitempty" json:"updatedAt"`
}

// CollectionName returns name of the models
func (a Admin) CollectionName() string {
	return "admins"
}

func initAdminIndex(db *mongo.Database) error {
	admin := Admin{}
	adminCol := db.Collection(admin.CollectionName())
	if err := createIndex(adminCol, bson.M{"phone": 1}, true); err != nil {
		return err
	}
	if err := createIndex(adminCol, bson.M{"email": 1}, true); err != nil {
		return err
	}
	return nil
}

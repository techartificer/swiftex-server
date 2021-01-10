package data

import (
	"context"
	"time"

	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/lib/jwt"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SessionRepository interface {
	CreateSession(db *mongo.Database, sess *models.Session) error
	UpdateSession(db *mongo.Database, token, accessToken string, userID primitive.ObjectID) (*models.Session, error)
	Logout(db *mongo.Database, token string) error
}

type SessionRepoImpl struct{}

var sessionRepo SessionRepository

func NewSessionRepo() SessionRepository {
	if sessionRepo == nil {
		sessionRepo = &SessionRepoImpl{}
	}
	return sessionRepo
}

func (s *SessionRepoImpl) CreateSession(db *mongo.Database, sess *models.Session) error {
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	_, err := sessionCollection.InsertOne(context.Background(), sess)
	return err
}

func (s *SessionRepoImpl) UpdateSession(db *mongo.Database, token, accessToken string, userID primitive.ObjectID) (*models.Session, error) {
	sess := &models.Session{}
	filter := bson.D{{"refreshToken", token}}
	after := options.After
	opt := options.FindOneAndUpdateOptions{
		ReturnDocument: &after,
	}
	update := bson.D{{"$set", bson.M{
		"refreshToken": jwt.NewRefresToken(userID),
		"accessToken":  accessToken,
		"createdAt":    time.Now().UTC(),
		"expiresOn":    time.Now().Add(time.Minute * time.Duration(config.GetJWT().RefreshTTL)),
	}}}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	err := sessionCollection.FindOneAndUpdate(context.Background(), filter, update, &opt).Decode(sess)
	return sess, err
}

func (s *SessionRepoImpl) Logout(db *mongo.Database, token string) error {
	sess := &models.Session{}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	filter := bson.D{{"refreshToken", token}}
	err := sessionCollection.FindOneAndDelete(context.Background(), filter).Decode(&sess)
	return err
}

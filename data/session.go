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
)

type SessionRepository interface {
	CreateSession(db *mongo.Database, sess *models.Session) error
	UpdateSession(db *mongo.Database, token, accessToken string, userID primitive.ObjectID) (*models.Session, error)
	Logout(db *mongo.Database, token string) error
	RemoveSessionsByUserID(db *mongo.Database, userID string) (*mongo.DeleteResult, error)
}

type sessionRepoImpl struct{}

var sessionRepo SessionRepository

func NewSessionRepo() SessionRepository {
	if sessionRepo == nil {
		sessionRepo = &sessionRepoImpl{}
	}
	return sessionRepo
}

func (s *sessionRepoImpl) CreateSession(db *mongo.Database, sess *models.Session) error {
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	_, err := sessionCollection.InsertOne(context.Background(), sess)
	return err
}

func (s *sessionRepoImpl) UpdateSession(db *mongo.Database, token, accessToken string, userID primitive.ObjectID) (*models.Session, error) {
	newsess := &models.Session{
		ID:           primitive.NewObjectID(),
		RefreshToken: jwt.NewRefresToken(userID),
		UserID:       userID,
		AccessToken:  accessToken,
		CreatedAt:    time.Now().UTC(),
		ExpiresOn:    time.Now().Add(time.Minute * time.Duration(config.GetJWT().RefreshTTL)),
	}
	filter := bson.D{{"refreshToken", token}}
	update := bson.D{{"$set", bson.M{"expiresOn": time.Now().Add(time.Second * 20)}}}
	collectionName := newsess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	sess := &models.Session{}
	err := sessionCollection.FindOneAndUpdate(context.Background(), filter, update).Decode(sess)
	if err != nil {
		return nil, err
	}
	_, err = sessionCollection.InsertOne(context.Background(), newsess)
	return newsess, err
}

func (s *sessionRepoImpl) Logout(db *mongo.Database, token string) error {
	sess := &models.Session{}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	filter := bson.D{{"refreshToken", token}}
	err := sessionCollection.FindOneAndDelete(context.Background(), filter).Decode(&sess)
	return err
}

func (s *sessionRepoImpl) RemoveSessionsByUserID(db *mongo.Database, userID string) (*mongo.DeleteResult, error) {
	_userID, err := primitive.ObjectIDFromHex(userID)
	if err != nil {
		return nil, err
	}
	sess := &models.Session{}
	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	filter := bson.D{{"userId", _userID}}
	res, err := sessionCollection.DeleteMany(context.Background(), filter)
	return res, err
}

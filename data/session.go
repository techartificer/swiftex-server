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
	sess := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       userID,
		RefreshToken: jwt.NewRefresToken(userID),
		AccessToken:  accessToken,
		CreatedAt:    time.Now().UTC(),
		ExpiresOn:    time.Now().Add(time.Minute * time.Duration(config.GetJWT().RefreshTTL)),
	}
	filter := bson.D{{"refreshToken", token}}
	update := bson.D{{"$set", bson.M{
		"expiresOn": time.Now().Add(time.Second * 5),
	}}}

	collectionName := sess.CollectionName()
	sessionCollection := db.Collection(collectionName)
	updateCh := make(chan error)
	insertCh := make(chan error)

	go func() {
		_, err := sessionCollection.UpdateOne(context.Background(), filter, update)
		updateCh <- err
	}()
	go func() {
		_, err := sessionCollection.InsertOne(context.Background(), sess)
		insertCh <- err
	}()

	if err := <-updateCh; err != nil {
		return nil, err
	}
	if err := <-insertCh; err != nil {
		return nil, err
	}
	return sess, nil
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

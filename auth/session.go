package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/AkifhanIlgaz/auth-mongo/rand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BytesPerToken = 32

type Session struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    string             `json:"userId"`
	SessionId string             `json:"sessionId"`
	Token     string             `json:"-"`
	TokenHash string             `json:"tokenHash"`
}

type SessionService struct {
	collection *mongo.Collection
}

func newSessionService(client *mongo.Client, database, collection string) *SessionService {
	return &SessionService{
		collection: client.Database(database).Collection(collection),
	}
}

func (service *SessionService) Create(userId string) (*Session, error) {
	token, err := rand.String(BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create session | rand: %w", err)
	}

	id := primitive.NewObjectID()
	session := Session{
		Id:        id,
		SessionId: id.Hex(),
		UserId:    userId,
		TokenHash: service.hash(token),
	}

	count, err := service.collection.CountDocuments(context.TODO(), bson.M{
		"userid": userId,
	})
	if err != nil {
		return nil, fmt.Errorf("create user & count: %w", err)
	}
	if count > 0 {
		// update
		res, err := service.collection.UpdateOne(context.TODO(), bson.M{
			"userid": userId,
		},
			bson.M{
				"$set": bson.M{
					"tokenhash": session.TokenHash,
				},
			},
		)
		if err != nil || res.ModifiedCount == 0 {
			return nil, fmt.Errorf("create session | update: %w", err)
		}
	} else {
		// insert
		_, err := service.collection.InsertOne(context.Background(), session)
		if err != nil {
			return nil, fmt.Errorf("create session | insert: %w", err)
		}
	}

	session.Token = token
	return &session, nil

}

func (service *SessionService) UserId(token string) (string, error) {
	tokenHash := service.hash(token)
	var session Session

	res := service.collection.FindOne(context.TODO(), bson.M{
		"tokenhash": tokenHash,
	})
	err := res.Decode(&session)
	if err != nil {
		return "", fmt.Errorf("session user: %w", err)
	}

	return session.UserId, nil
}

func (service *SessionService) Delete(token string) error {
	tokenHash := service.hash(token)

	res, err := service.collection.DeleteOne(context.TODO(), bson.M{
		"tokenhash": tokenHash,
	})
	if err != nil {
		return fmt.Errorf("delete session: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("Cannot find session for given token")
	}

	return nil
}

func (service *SessionService) hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

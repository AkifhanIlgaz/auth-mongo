package auth

import (
	"crypto/sha256"
	"encoding/base64"
	"fmt"

	"github.com/AkifhanIlgaz/auth-mongo/rand"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const BytesPerToken = 32

type Session struct {
	Id        primitive.ObjectID `bson:"_id"`
	UserId    string             `json:"userId"`
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
		UserId:    id.Hex(),
		Token:     token,
		TokenHash: service.hash(token),
	}

	// Check if user has valid session token
	// If so, update session token with the new one

	return &session, nil

}

func (service *SessionService) hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

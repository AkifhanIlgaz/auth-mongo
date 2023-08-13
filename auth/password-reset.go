package auth

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/AkifhanIlgaz/auth-mongo/rand"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const TimeDuration = 1 * time.Hour

type PasswordResetService struct {
	collection     *mongo.Collection
	userCollection *mongo.Collection
}

type PasswordReset struct {
	Id              primitive.ObjectID `bson:"_id`
	PasswordResetId string             `json:"passwordResetId`
	UserId          string             `json:"userId"`
	Token           string             `json"-"`
	TokenHash       string             `json:"tokenHash"`
	ExpiresAt       time.Time          `json:"expiresAt"`
}

func newPasswordResetService(client *mongo.Client, database, collection string) *PasswordResetService {
	return &PasswordResetService{
		collection:     client.Database(database).Collection(collection),
		userCollection: client.Database(database).Collection("Users"),
	}
}

func (service *PasswordResetService) Create(email string) (*PasswordReset, error) {
	email = strings.TrimSpace(email)

	var user User
	err := service.userCollection.FindOne(context.TODO(), bson.M{
		"email": email,
	}).Decode(&user)
	if err != nil {
		return nil, fmt.Errorf("create password reset | find user: %w", err)
	}
	if user.UserId == "" {
		return nil, fmt.Errorf("user with this email doesn't exist")
	}

	token, err := rand.String(BytesPerToken)
	if err != nil {
		return nil, fmt.Errorf("create password reset | rand string: %w", err)
	}

	id := primitive.NewObjectID()
	passwordReset := PasswordReset{
		Id:              id,
		PasswordResetId: id.Hex(),
		UserId:          user.UserId,
		TokenHash:       service.hash(token),
		ExpiresAt:       time.Now().Add(TimeDuration),
	}

	count, err := service.collection.CountDocuments(context.TODO(), bson.M{
		"userid": user.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("create password reset & count: %w", err)
	}
	if count > 0 {
		res, err := service.collection.UpdateOne(context.TODO(), bson.M{
			"userid": user.UserId,
		},
			bson.M{
				"$set": bson.M{
					"tokenhash": passwordReset.TokenHash,
					"expiresat": passwordReset.ExpiresAt,
				},
			},
		)
		if err != nil || res.ModifiedCount == 0 {
			return nil, fmt.Errorf("create password reset | update: %w", err)
		}
	} else {
		_, err := service.collection.InsertOne(context.Background(), passwordReset)
		if err != nil {
			return nil, fmt.Errorf("create password reset | insert: %w", err)
		}
	}

	passwordReset.Token = token
	return &passwordReset, nil
}

func (service *PasswordResetService) Consume(token string) error {
	tokenHash := service.hash(token)
	var passwordReset PasswordReset

	err := service.collection.FindOne(context.TODO(), bson.M{
		"tokenhash": tokenHash,
	}).Decode(&passwordReset)
	if err != nil {
		return fmt.Errorf("consume password reset token: %w", err)
	}

	if time.Now().After(passwordReset.ExpiresAt) {
		return fmt.Errorf("token expired: %v", token)
	}

	err = service.delete(passwordReset.PasswordResetId)
	if err != nil {
		return fmt.Errorf("consume | delete: %w", err)
	}

	return nil
}

func (service *PasswordResetService) delete(id string) error {
	res, err := service.collection.DeleteOne(context.TODO(), bson.M{
		"passwordresetid": id,
	})
	if err != nil {
		return fmt.Errorf("delete password reset: %w", err)
	}
	if res.DeletedCount == 0 {
		return fmt.Errorf("there is no password reset for this id")
	}

	return nil
}

func (service *PasswordResetService) hash(token string) string {
	hash := sha256.Sum256([]byte(token))
	return base64.URLEncoding.EncodeToString(hash[:])
}

package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           primitive.ObjectID `bson:"_id"`
	Email        string             `json:"email"`
	PasswordHash string             `json:"passwordHash"`
	CreatedAt    time.Time          `json:"createdAt"`
	UserId       string             `json:"userId"`
}

type UserService struct {
	collection *mongo.Collection
}

func newUserService(client *mongo.Client, database, collection string) *UserService {
	return &UserService{
		collection: client.Database(database).Collection(collection),
	}
}

// Make sure that email is valid and password is strong on front-end
func (service *UserService) Create(email, password string) (*User, error) {
	email = strings.TrimSpace(email)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	id := primitive.NewObjectID()
	user := User{
		Id:           id,
		Email:        email,
		PasswordHash: string(passwordHash),
		CreatedAt:    time.Now(),
		UserId:       id.Hex(),
	}

	count, err := service.collection.CountDocuments(context.TODO(), bson.M{
		"email": email,
	})
	if err != nil {
		return nil, fmt.Errorf("create user & count: %w", err)
	}
	if count > 0 {
		return nil, fmt.Errorf("email already exists")
	}

	_, err = service.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

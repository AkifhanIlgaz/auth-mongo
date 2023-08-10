package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

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
	Collection *mongo.Collection
}

func newUserService(client *mongo.Client, database, collection string) *UserService {
	return &UserService{
		Collection: client.Database(database).Collection(collection),
	}
}

func (service *UserService) Create(email, password string) (*User, error) {
	// Make sure this is a valid email on front-end
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
		UserId:       id.String(),
	}

	// TODO: Insert new user to Mongo
	// ? Mongo DB constraints
	res, err := service.Collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	// ! Email must be unique to user

	return nil, nil
}

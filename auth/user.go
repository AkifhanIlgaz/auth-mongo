package auth

import (
	"fmt"
	"strings"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Id           string // TODO: Mongo object id
	Email        string
	PasswordHash string
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

	// TODO: Insert new user to Mongo
	//
}

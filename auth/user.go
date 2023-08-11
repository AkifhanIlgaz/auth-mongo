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

// ! Make sure that email is valid and password is strong on front-end
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
		return nil, ErrEmailTaken
	}

	_, err = service.collection.InsertOne(context.TODO(), user)
	if err != nil {
		return nil, fmt.Errorf("create user: %w", err)
	}

	return &user, nil
}

func (service *UserService) Delete(userId string) error {
	res, err := service.collection.DeleteOne(context.TODO(), bson.M{
		"userid": userId,
	})
	if err != nil {
		return fmt.Errorf("delete user: %w", err)
	}
	if res.DeletedCount == 0 {
		return ErrUserDoesntExist
	}

	return nil
}

func (service *UserService) UpdatePassword(userId string, password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("update password | bcrypt: %w", err)
	}

	res, err := service.collection.UpdateOne(
		context.TODO(), bson.M{
			"userid": userId,
		}, bson.M{
			"$set": bson.M{
				"passwordhash": string(passwordHash),
			},
		},
	)
	if err != nil {
		return fmt.Errorf("update password | update one: %w", err)
	}
	if res.ModifiedCount == 0 {
		return ErrUserDoesntExist
	}

	return nil
}

func (service *UserService) GetUser(userId string) User {
	var user User

	res := service.collection.FindOne(context.TODO(), bson.M{
		"userid": userId,
	})

	res.Decode(&user)

	return user
}

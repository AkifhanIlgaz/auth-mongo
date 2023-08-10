package auth

import "go.mongodb.org/mongo-driver/mongo"

type UserService struct {
	DB *mongo.Database
}

func newUserService(client *mongo.Client) *UserService {
	return &UserService{
		DB: client.Database("users"),
	}
}

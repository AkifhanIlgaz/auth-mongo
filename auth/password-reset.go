package auth

import "go.mongodb.org/mongo-driver/mongo"

type PasswordResetService struct {
	DB *mongo.Database
}

func newPasswordResetService(client *mongo.Client) *PasswordResetService {
	return &PasswordResetService{
		DB: client.Database("password-reset-tokens"),
	}
}

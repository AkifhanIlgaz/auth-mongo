package auth

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AuthService struct {
	UserService          *UserService
	SessionService       *SessionService
	PasswordResetService *PasswordResetService
}

func New(mongoUri string) (*AuthService, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(mongoUri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, fmt.Errorf("new auth : %w", err)
	}
	

	if err := client.Database("admin").RunCommand(context.TODO(), bson.D{{"ping", 1}}).Err(); err != nil {
		return nil, fmt.Errorf("new auth: %w", err)
	}
	fmt.Println("Pinged your deployment. You successfully connected to MongoDB!")

	return &AuthService{
		UserService:          newUserService(client),
		SessionService:       newSessionService(client),
		PasswordResetService: newPasswordResetService(client),
	}, nil
}

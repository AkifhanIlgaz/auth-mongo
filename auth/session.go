package auth

import "go.mongodb.org/mongo-driver/mongo"

type SessionService struct {
	DB *mongo.Database
}

func newSessionService(client *mongo.Client) *SessionService {
	return &SessionService{
		DB: client.Database("sessions"),
	}
}

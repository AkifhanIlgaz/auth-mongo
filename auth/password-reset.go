package auth

import "go.mongodb.org/mongo-driver/mongo"

type PasswordResetService struct {
	Collection *mongo.Collection
}

func newPasswordResetService(client *mongo.Client, database, collection string) *PasswordResetService {
	return &PasswordResetService{
		Collection: client.Database(database).Collection(collection),
	}
}

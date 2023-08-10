package auth

import "go.mongodb.org/mongo-driver/mongo"

type SessionService struct {
	Collection *mongo.Collection
}

func newSessionService(client *mongo.Client, database, collection string) *SessionService {
	return &SessionService{
		Collection: client.Database(database).Collection(collection),
	}
}

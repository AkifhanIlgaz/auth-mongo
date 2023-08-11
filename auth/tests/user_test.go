package auth

import (
	"os"
	"testing"

	"github.com/AkifhanIlgaz/auth-mongo/auth"
	"github.com/joho/godotenv"
)

var authService *auth.AuthService

func init() {
	var err error
	godotenv.Load("../../.env")
	authService, err = auth.New(os.Getenv("MONGO_URI"), "VocabBuilder")
	if err != nil {
		panic(err)
	}
}

func TestCreateUser(t *testing.T) {
	email, password := "test@gmail.com", "testing"

	// Create new user
	if _, err := authService.UserService.Create(email, password); err != nil {
		if err != auth.ErrEmailTaken {
			t.Fatalf("cannot create user, %v", err)
		}
	}

	// Cannot create user with an existing email
	if _, err := authService.UserService.Create(email, password); err != auth.ErrEmailTaken {
		t.Fatal("creating user with an existing email")
	}
}

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

}

func TestDuplicateUser(t *testing.T) {
	email, password := "testduplicate@gmail.com", "testing"

	// Create new user
	if _, err := authService.UserService.Create(email, password); err != nil {
		if err != auth.ErrEmailTaken {
			t.Fatalf("cannot create user, %v", err)
		}
	}

	// Create another user with same email
	if _, err := authService.UserService.Create(email, password); err == nil {
		t.Fatalf("creating user with duplicate email")
	}
}

func TestDeleteUser(t *testing.T) {
	email, password := "testdelete@gmail.com", "testing"

	user, err := authService.UserService.Create(email, password)
	if err != nil {
		if err != auth.ErrEmailTaken {
			t.Fatalf("cannot create user, %v", err)
		}
	}

	// Delete user
	authService.UserService.Delete(user.UserId)

	// Check if it is successfully deleted
	if err := authService.UserService.Delete(user.UserId); err == nil {
		t.Fatal("non-existing user is deleted")
	}
}

func TestNonExistingUser(t *testing.T) {
	userid := "1"

	err := authService.UserService.Delete(userid)
	if err == nil {
		t.Fatal("deleting non-existing user")
	}
}

func TestUpdatePassword(t *testing.T) {
	email, password := "testupdatepassword@gmail.com", "testing"

	user, err := authService.UserService.Create(email, password)
	if err != nil {
		if err != auth.ErrEmailTaken {
			t.Fatalf("cannot create user, %v", err)
		}
	}

	if err = authService.UserService.UpdatePassword(user.UserId, "nesdşflksf"); err != nil {
		t.Fatal("cannot update password")
	}

	newUser := authService.UserService.GetUser(user.UserId)

	if newUser.PasswordHash == user.PasswordHash {
		t.Fatal("password hashes are the same")
	}
}

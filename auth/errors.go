package auth

import "errors"

var (
	ErrEmailTaken      error = errors.New("Email is taken by another user")
	ErrUserDoesntExist error = errors.New("User doesn't exist")
)

package auth

import "errors"

var (
	ErrEmailTaken error = errors.New("Email is taken by another user")
)

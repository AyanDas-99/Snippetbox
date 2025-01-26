package models

import (
	"errors"
)

var (
	ErrNoRecord = errors.New("models: no matching record found")

	// if user logs in with incorrect email or password
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// if user tries signup with existing email
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

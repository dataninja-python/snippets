package models

import "errors"

var (
	ErrNoRecord = errors.New("models: no matching record found")
	// ErrInvalidCredentials deals with incorrect emails or passwords
	ErrInvalidCredentials = errors.New("models: invalid credentials")
	// ErrDuplicateEmail deals with an email already existing in the database
	ErrDuplicateEmail = errors.New("models: duplicate email")
)

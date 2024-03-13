package models

import (
	"database/sql"
	"time"
)

// User struct represents users data in the application
type User struct {
	ID             int
	Name           string
	Email          string
	HashedPassword byte
	Created        time.Time
}

// UserModel wraps the database connection pool
type UserModel struct {
	DB *sql.DB
}

// Insert adds a new record to the users table
func (m *UserModel) Insert(name, email, password string) error {
	return nil
}

// Authenticate verifies if a user exists using an email and password. If so, it returns the user ID
func (m *UserModel) Authenticate(email, password string) (int, error) {
	return 0, nil
}

// Exists checks if a user exists with a specific ID
func (m *UserModel) Exists(id int) (bool, error) {
	return false, nil
}

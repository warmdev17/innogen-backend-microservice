package models

import "time"

// User represents a row in the users table.
// Password is tagged json:"-" to prevent accidental serialization in API responses.
type User struct {
	ID        int       `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`
	Username  *string   `json:"username"`
	FullName  *string   `json:"fullName"`
	Role      string    `json:"role"`
	IsActive  bool      `json:"isActive"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

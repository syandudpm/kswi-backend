package person

import "time"

// CreatePersonRequest represents the request payload for creating a person
type CreatePersonRequest struct {
	PersonID string `json:"person_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PersonResponse represents the response structure for person data
type PersonResponse struct {
	ID        uint      `json:"id"`
	PersonID  string    `json:"person_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

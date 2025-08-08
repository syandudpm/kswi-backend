package people

import "time"

// CreatePersonRequest represents the request payload for creating a people
type CreatePersonRequest struct {
	PersonID string `json:"person_id" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// PersonResponse represents the response structure for people data
type PersonResponse struct {
	ID        uint      `json:"id"`
	PersonID  string    `json:"people_id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

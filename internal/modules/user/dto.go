package user

type CreateUserRequest struct {
	PersonID string `json:"person_id" binding:"required,min=1,max=150"`
	Username string `json:"username" binding:"required,min=3,max=150,alphanum"`
	Password string `json:"password" binding:"required,min=6"`
}

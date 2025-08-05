package user

import (
	"errors"
	"kswi-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists = errors.New("username already exists")
	ErrInvalidInput      = errors.New("invalid input")
)

type Service interface {
	CreateUser(req *CreateUserRequest) (*model.User, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func (s *service) CreateUser(req *CreateUserRequest) (*model.User, error) {
	// Optional: validate input
	if req.Username == "" || req.Password == "" || req.PersonID == "" {
		return nil, ErrInvalidInput
	}

	// Check if username already exists
	_, err := s.repo.FindByUsername(req.Username)
	if err == nil {
		return nil, ErrUserAlreadyExists
	}

	hashedPassword, err := hashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		PersonID: req.PersonID,
		Username: req.Username,
		Password: hashedPassword, // Store hashed password
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, err
	}

	// Clear password for safety (though it's already excluded in JSON via "-")
	user.Password = ""

	return user, nil
}

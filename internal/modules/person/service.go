package person

import (
	"errors"
	"kswi-backend/internal/model"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	CreatePerson(req *CreatePersonRequest) (*PersonResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) CreatePerson(req *CreatePersonRequest) (*PersonResponse, error) {
	// Validate required fields
	if req.PersonID == "" {
		return nil, errors.New("person_id is required")
	}
	if req.Username == "" {
		return nil, errors.New("username is required")
	}
	if req.Password == "" {
		return nil, errors.New("password is required")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create person model
	person := &model.Person{
		PersonID: req.PersonID,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	// Save to database
	err = s.repo.Create(person)
	if err != nil {
		return nil, err
	}

	// Return response
	response := &PersonResponse{
		ID:        person.ID,
		PersonID:  person.PersonID,
		Username:  person.Username,
		CreatedAt: person.CreatedAt,
		UpdatedAt: person.UpdatedAt,
	}

	return response, nil
}

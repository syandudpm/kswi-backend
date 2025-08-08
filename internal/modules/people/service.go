package people

import (
	"kswi-backend/internal/model"
	"kswi-backend/internal/shared/errors"

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
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.NewInternalError(err)
	}

	// Create people model
	people := &model.People{
		PersonID: req.PersonID,
		Username: req.Username,
		Password: string(hashedPassword),
	}

	// Save to database
	err = s.repo.Create(people)
	if err != nil {
		return nil, errors.NewDatabaseError(err)
	}

	// Return response
	response := &PersonResponse{
		ID:        people.ID,
		PersonID:  people.PersonID,
		Username:  people.Username,
		CreatedAt: people.CreatedAt,
		UpdatedAt: people.UpdatedAt,
	}

	return response, nil
}

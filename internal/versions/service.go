package versions

import (
	"time"

	"github.com/google/uuid"
)

type Service interface {
	Create(request *Request) error
	GetAll() ([]Version, error)
	GetById(id string) (*Version, error)
	Update(id string, request *Request) error
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(request *Request) error {
	postId, err := uuid.Parse(request.PostId)
	if err != nil {
		return err
	}
	
	version := Version{
		ID:           uuid.New(),
		Version:      request.Version,
		PostId:       postId,
		VersionNumber: request.VersionNumber,
		Content:      request.Content,
		CreatedAt:    time.Now(),
	}
	return s.repo.Create(&version)
}

func (s *service) GetAll() ([]Version, error) {
	return s.repo.GetAll()
}

func (s *service) GetById(id string) (*Version, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repo.GetById(parsedId)
}

func (s *service) Update(id string, request *Request) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	parsedPostId, err := uuid.Parse(request.PostId)
	if err != nil {
		return err
	}
	version := Version{
		ID:           parsedId,
		Version:      request.Version,
		PostId:       parsedPostId,
		VersionNumber: request.VersionNumber,
		Content:      request.Content,
		CreatedAt:    time.Now(),
	}
	return s.repo.Update(parsedId, &version)
}

func (s *service) Delete(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(parsedId)
}

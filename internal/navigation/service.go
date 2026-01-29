package navigation

import "github.com/google/uuid"

type Service interface {
	Create(navigation *Request) (*Navigation, error)
	GetAll() ([]Navigation, error)
	GetById(id string) (*Navigation, error)
	GetBySlug(slug string) (*Navigation, error)
	Update(id string, navigation *Request) (*Navigation, error)
	Delete(id string) error
}

type service struct {
	repository Repository
}

func NewService(repository Repository) Service {
	return &service{repository: repository}
}

func (s *service) Create(request *Request) (*Navigation, error) {
	var parentID *uuid.UUID
	if request.ParentId != "" {
		parsed, err := uuid.Parse(request.ParentId)
		if err != nil {
			return nil, err
		}
		parentID = &parsed
	}

	var categoryID *uuid.UUID
	if request.CategoryId != "" {
		parsed, err := uuid.Parse(request.CategoryId)
		if err != nil {
			return nil, err
		}
		categoryID = &parsed
	}

	var postID *uuid.UUID
	if request.PostId != "" {
		parsed, err := uuid.Parse(request.PostId)
		if err != nil {
			return nil, err
		}
		postID = &parsed
	}

	navigation := Navigation{
		ID:         uuid.New(),
		Label:      request.Label,
		Slug:       request.Slug,
		Type:       request.Type,
		Order:      request.Order,
		ParentId:   parentID,
		LinkSource: request.LinkSource,
		CategoryId: categoryID,
		PostId:     postID,
	}
	if err := s.repository.Create(&navigation); err != nil {
		return nil, err
	}
	return &navigation, nil
}

func (s *service) GetAll() ([]Navigation, error) {
	return s.repository.GetAll()
}

func (s *service) GetById(id string) (*Navigation, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return s.repository.GetById(parsedId)
}

func (s *service) GetBySlug(slug string) (*Navigation, error) {
	return s.repository.GetBySlug(slug)
}

func (s *service) Update(id string, request *Request) (*Navigation, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	var parentID *uuid.UUID
	if request.ParentId != "" {
		parsed, err := uuid.Parse(request.ParentId)
		if err != nil {
			return nil, err
		}
		parentID = &parsed
	}

	var categoryID *uuid.UUID
	if request.CategoryId != "" {
		parsed, err := uuid.Parse(request.CategoryId)
		if err != nil {
			return nil, err
		}
		categoryID = &parsed
	}

	var postID *uuid.UUID
	if request.PostId != "" {
		parsed, err := uuid.Parse(request.PostId)
		if err != nil {
			return nil, err
		}
		postID = &parsed
	}

	navigation := Navigation{
		ID:         parsedId,
		Label:      request.Label,
		Slug:       request.Slug,
		Type:       request.Type,
		Order:      request.Order,
		ParentId:   parentID,
		LinkSource: request.LinkSource,
		CategoryId: categoryID,
		PostId:     postID,
	}
	if err := s.repository.Update(parsedId, &navigation); err != nil {
		return nil, err
	}
	return &navigation, nil
}

func (s *service) Delete(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repository.Delete(parsedId)
}

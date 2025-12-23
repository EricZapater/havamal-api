package categories

import "github.com/google/uuid"

type Service interface {
	Create(category *Request) error
	GetAll() ([]Category, error)
	GetById(id string) (*Category, error)
	GetBySlug(slug string) (*Category, error)
	Update(id string, category *Request) error
	Delete(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (service *service) Create(category *Request) error {
	return service.repo.Create(&Category{
		ID:          uuid.New(),
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Order:       category.Order,
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	})
}

func (service *service) GetAll() ([]Category, error) {
	return service.repo.GetAll()
}

func (service *service) GetById(id string) (*Category, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return service.repo.GetById(parsedId)
}

func (service *service) GetBySlug(slug string) (*Category, error) {
	return service.repo.GetBySlug(slug)
}

func (service *service) Update(id string, category *Request) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return service.repo.Update(&Category{
		ID:          parsedId,
		Name:        category.Name,
		Slug:        category.Slug,
		Description: category.Description,
		Order:       category.Order,
		UpdatedAt:   category.UpdatedAt,
	})
}

func (service *service) Delete(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return service.repo.Delete(parsedId)
}
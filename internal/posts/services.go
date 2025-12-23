package posts

import (
	"context"
	"havamal-api/internal/users"

	"github.com/google/uuid"
)

type Service interface {
	CreatePost(post *Request) error
	GetPost(id string) (*Response, error)
	GetPosts() ([]Response, error)
	GetPublishedPosts() ([]Response, error) 
	GetPostsByAuthor(authorId string) ([]Response, error)
	GetPostBySlug(slug string) (*Response, error)
	GetSummariesByCategory(category string) ([]Response, error)
	UpdatePost(id string, post *Request) error
	DeletePost(id string) error
	AddCategory(request PostCategories) error
	DeleteCategory(request PostCategories) error
	AddVersion(request PostVersion) error
	DeleteVersion(request PostVersion) error
}

type service struct {
	repo        Repository
	userService users.Service
}

func NewService(repo Repository, userService users.Service) Service {
	return &service{
		repo:        repo,
		userService: userService,
	}
}

func (service *service) CreatePost(post *Request) error {
	var authorId uuid.UUID
	var err error

	if post.AuthorId == "" && post.Author != "" {
		// Look up user by email
		user, err := service.userService.FindByEmail(context.Background(), post.Author)
		if err != nil {
			return err
		}
		authorId = user.ID
	} else {
		authorId, err = uuid.Parse(post.AuthorId)
		if err != nil {
			return err
		}
	}

	newPostId := uuid.New()
	err = service.repo.CreatePost(&Post{
		ID:          newPostId,
		Title:       post.Title,
		Slug:        post.Slug,
		Summary:     post.Summary,
		Content:     post.Content,
		Status:      post.Status,
		PublishedAt: post.PublishedAt,
		UpdatedAt:   post.UpdatedAt,
		AuthorId:    authorId,
	})
	if err != nil {
		return err
	}

	if post.CategoryId != "" {
		categoryId, err := uuid.Parse(post.CategoryId)
		if err == nil {
			_ = service.repo.AddCategory(PostCategories{
				PostId:     newPostId,
				CategoryId: categoryId,
			})
		}
	}
	return nil
}

func (service *service) GetPost(id string) (*Response, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	return service.repo.GetPost(parsedId)
}

func (service *service) GetPosts() ([]Response, error) {
	return service.repo.GetPosts()
}

func (service *service) GetPublishedPosts() ([]Response, error) {
	return service.repo.GetPublishedPosts()
}

func (service *service) GetPostsByAuthor(authorId string) ([]Response, error) {
	parsedAuthorId, err := uuid.Parse(authorId)
	if err != nil {
		return nil, err
	}
	return service.repo.GetPostsByAuthor(parsedAuthorId)
}

func (service *service) GetPostBySlug(slug string) (*Response, error) {
	return service.repo.GetPostBySlug(slug)
}

func (service *service) GetSummariesByCategory(category string) ([]Response, error) {
	return service.repo.GetSummariesByCategory(category)
}

func (service *service) UpdatePost(id string, post *Request) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	// Fetch existing post to preserve data like AuthorId if not provided
	existingPost, err := service.repo.GetPost(parsedId)
	if err != nil {
		return err
	}

	authorId := existingPost.AuthorId
	if post.AuthorId != "" {
		parsedAuthorId, err := uuid.Parse(post.AuthorId)
		if err == nil {
			authorId = parsedAuthorId
		}
	}

	return service.repo.UpdatePost(&Post{
		ID:          parsedId,
		Title:       post.Title,
		Slug:        post.Slug,
		Summary:     post.Summary,
		Content:     post.Content,
		Status:      post.Status,
		PublishedAt: post.PublishedAt,
		UpdatedAt:   post.UpdatedAt,
		AuthorId:    authorId,
	})
}

func (service *service) DeletePost(id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return service.repo.DeletePost(parsedId)
}

func (service *service) AddCategory(request PostCategories) error {
	return service.repo.AddCategory(request)
}

func (service *service) DeleteCategory(request PostCategories) error {
	return service.repo.DeleteCategory(request)
}

func (service *service) AddVersion(request PostVersion) error {
	return service.repo.AddVersion(request)
}

func (service *service) DeleteVersion(request PostVersion) error {
	return service.repo.DeleteVersion(request)
}

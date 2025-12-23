package users

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	Create(ctx context.Context, request UserRequest) (User, error)	
	FindAll(ctx context.Context) ([]User, error)	
	FindByID(ctx context.Context, id string) (User, error)
	FindByEmail(ctx context.Context, email string) (User, error)
	Update(ctx context.Context, id string, request UserRequest) (User, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, request UserRequest) (User, error) {
	
	// Recuperar i validar is_admin
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return User{}, errors.New("invalid or missing is_admin in context")
	}
	if !isAdmin{
		return User{}, errors.New("user is not admin")
	}
	
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}
	
	user := User{
		ID:         uuid.New(),
		Username:   request.Username,
		Email:      request.Email,
		Password:   string(hashedPassword),
		IsActive:   true,
		IsAdmin:    false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	return s.repo.Create(ctx, user)
}


func (s *service) FindAll(ctx context.Context) ([]User, error) {
	isAdminVal := ctx.Value("is_admin")
	isAdmin, ok := isAdminVal.(bool)
	if !ok {
		return nil, errors.New("invalid or missing is_admin in context")
	}
	if isAdmin{
		return s.repo.FindAll(ctx)
	}
	return nil, errors.New("user is not admin")
}

func (s *service) FindByID(ctx context.Context, id string) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}
	return s.repo.FindByID(ctx, parsedId)
}

func (s *service) FindByEmail(ctx context.Context, email string) (User, error) {
	return s.repo.FindByEmail(ctx, email)
}

func (s *service) Update(ctx context.Context, id string, request UserRequest) (User, error) {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return User{}, err
	}

	user, err := s.repo.FindByID(ctx, parsedId)
	if err != nil {
		return User{}, err
	}
	
	user.Username = request.Username
	user.Email = request.Email
	
	// Only hash and update password if a new one is provided
	if request.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
		if err != nil {
			return User{}, err
		}
		user.Password = string(hashedPassword)
	}

	user.IsActive = request.IsActive
	user.IsAdmin = false
	user.UpdatedAt = time.Now()
	return s.repo.Update(ctx, user)
}

func (s *service) Delete(ctx context.Context, id string) error {
	parsedId, err := uuid.Parse(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(ctx, parsedId)
}

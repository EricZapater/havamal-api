package auth

import (
	"context"
	"havamal-api/internal/users"
	"havamal-api/middleware"

	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
    Login(ctx context.Context, req LoginRequest) (string, users.User, time.Time, error)
    Register(ctx context.Context, req RegisterRequest) (users.User, error)
    ValidateUser(ctx context.Context, username, password string) (users.User, error)
}

type authService struct {
	userService users.Service
	jwtMiddleware *jwt.GinJWTMiddleware
}

func NewAuthService(userService users.Service, jwtMiddleware *jwt.GinJWTMiddleware) AuthService {
	return &authService{
		userService: userService,
		jwtMiddleware: jwtMiddleware,
	}
}

// Login verifica les credencials i retorna un token JWT si són vàlides
func (s *authService) Login(ctx context.Context, req LoginRequest) (string, users.User, time.Time, error) {
    // Validar les credencials
    user, err := s.ValidateUser(ctx, req.Email, req.Password)
    if err != nil {
        return "", users.User{}, time.Time{}, err
    }
    
    // Generar token JWT
    authUser := &middleware.AuthUser{
        ID:         user.ID.String(),
        Username:   user.Username,
        Email:      user.Email,        
        IsAdmin:    user.IsAdmin,
    }
    token, expire, err := s.jwtMiddleware.TokenGenerator(authUser)
    if err != nil {
        return "", users.User{}, time.Time{}, err
    }
   
    user.Password = "" // No retornar la contrasenya en la resposta
    return token, user,  expire, nil
}

func (s *authService) Register(ctx context.Context, req RegisterRequest) (users.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return users.User{}, err
    }

    user := users.UserRequest{
        Email:      req.Email,
        Password:   string(hashedPassword),
    }

    return s.userService.Create(ctx, user)
}

// ValidateUser verifica si les credencials són vàlides i retorna l'ID de l'usuari
func (s *authService) ValidateUser(ctx context.Context, email, password string) (users.User, error) {
    // Obtenir l'usuari per email
    user, err := s.userService.FindByEmail(ctx, email)

    if err != nil {
        return users.User{}, ErrUserNotFound
    }
    
    // Verificar que l'usuari estigui actiu
    if !user.IsActive {
        return users.User{}, ErrInactiveUser
    }
    
    // Verificar la contrasenya
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
    if err != nil {
        return users.User{}, ErrInvalidCredentials
    }
    
    // Retornar l'ID de l'usuari com a identificador principal
    return user, nil
}
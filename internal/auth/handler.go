package auth

import (
	"net/http"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
    authService    AuthService
    jwtMiddleware *jwt.GinJWTMiddleware
}

func NewAuthHandler(authService AuthService, jwtMiddleware *jwt.GinJWTMiddleware) *AuthHandler {
    return &AuthHandler{
        authService:    authService,
        jwtMiddleware: jwtMiddleware,
    }
}

// Login processa una petició de login i retorna un token JWT
func (h *AuthHandler) Login(c *gin.Context) {
    var loginRequest LoginRequest
    if err := c.ShouldBindJSON(&loginRequest); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid login request"})
        return
    }
    
    token,user, expire, err := h.authService.Login(c.Request.Context(), loginRequest)
    if err != nil {
        var statusCode int
        switch err {
        case ErrInvalidCredentials:
            statusCode = http.StatusUnauthorized
        case ErrInactiveUser:
            statusCode = http.StatusForbidden
        default:
            statusCode = http.StatusInternalServerError
        }
        
        c.JSON(statusCode, gin.H{"error": err.Error()})
        return
    }
    
    c.JSON(http.StatusOK, LoginResponse{
        Token:  token,
        Expire: expire.Format(time.RFC3339),
        User:   user,
    })
}

// RefreshToken utilitza el middleware JWT per refrescar el token
func (h *AuthHandler) RefreshToken(c *gin.Context) {
    h.jwtMiddleware.RefreshHandler(c)
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.Register(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		ID:       user.ID.String(),
		Email:    user.Email,
	})
}

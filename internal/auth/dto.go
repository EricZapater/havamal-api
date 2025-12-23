package auth

import (
	"havamal-api/internal/users"
)

type LoginRequest struct {
	Email string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token  string `json:"token"`
	Expire string `json:"expire"`
	User   users.User   `json:"user"`	
}

type RegisterRequest struct {
	Password   string `json:"password" binding:"required"`
	Email      string `json:"email" binding:"required,email"`
	IsAdmin    bool   `json:"is_admin"`
}

type RegisterResponse struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
}
package handlers

import "ROP_Backend/internal/models"

type UserResponse struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  *models.User `json:"user"`
}

type OTPResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

package handlers

import "ROP_Backend/internal/models"

type UserResponse struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
}

type OnboardingResponse struct {
	Message string `json:"message"`
}

type LoginResponse struct {
	Token          string       `json:"token"`
	User           *models.User `json:"user"`
	NeedOnboarding bool         `json:"needOnboarding"`
}

type OTPResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type OTPErrorInfo struct {
	Code   string `json:"code"`
	Detail string `json:"detail,omitempty"`
}

type OTPResponse struct {
	Success bool          `json:"success"`
	Message string        `json:"message"`
	Data    interface{}   `json:"data,omitempty"`
	Error   *OTPErrorInfo `json:"error,omitempty"`
}

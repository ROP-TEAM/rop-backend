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

type ErrorResponse struct {
	Error string `json:"error"`
}

type APIResponse struct {
	Status  string      `json:"status" example:"success"` //suc/fail/err
	Message string      `json:"message" example:"OTP sent successfully"`
	Data    interface{} `json:"data,omitempty"`
}

func Success(message string, data interface{}) APIResponse {
	return APIResponse{
		Status:  "success",
		Message: message,
		Data:    data,
	}
}

func Fail(message string) APIResponse {
	return APIResponse{
		Status:  "fail",
		Message: message,
	}
}

func Error(message string) APIResponse {
	return APIResponse{
		Status:  "error",
		Message: message,
	}
}

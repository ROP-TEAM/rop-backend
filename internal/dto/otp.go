package dto

import "time"

type VerifyOTPRequest struct {
	Pin   string `json:"pin" validate:"required,len=10,numeric" example:"123456"`
	Tel   string `json:"tel" validate:"required" example:"0812345678"`
	RefNo string `json:"refNo" validate:"required" example:"ABC123"`
}

type RequestOTPRequest struct {
	Tel string `json:"tel" validate:"required,len=10,numeric" example:"0812345678"`
}

type RequestOTPResponse struct {
	RefNo       string    `json:"refNo" example:"ABC123"`
	Status      string    `json:"status" example:"success"`
	Tel         string    `json:"tel" example:"0999999999"`
	MaxAttempts int       `json:"max_attempts" example:"5"`
	ExpiresAt   time.Time `json:"expires_at"`
	Attempts    int       `json:"attempts" example:"0"`
}

type VerifyOTPResponse struct {
	Status string `json:"status" example:"success"`
}

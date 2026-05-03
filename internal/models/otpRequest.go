package models

import "time"

type OtpRequest struct {
	ID     uint   `gorm:"primaryKey"`
	Tel    string `gorm:"size:10;index;not null;"`
	RefNo  string `gorm:"size:50;index;not null;"`
	Token  string `gorm:"size:255; not null"`
	UserID uint   `gorm:"not null" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user"`

	ExpiresAt   time.Time `gorm:"index"`
	Attempts    int
	MaxAttempts int
	IsUsed      bool `gorm:"default:false"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type RequestOTPRequest struct {
	Tel string `json:"tel" example:"0812345678"`
}

type RequestOTPResponse struct {
	RefNo       string    `json:"refNo" example:"ABC123"`
	Status      string    `json:"status" example:"success"`
	Tel         string    `json:"tel" example:"0999999999"`
	MaxAttempts int       `json:"max_attempts" example:"5"`
	ExpiresAt   time.Time `json:"expires_at"`
	Attempts    int       `json:"attempts" example:"0"`
}

type VerifyOTPRequest struct {
	Pin   string `json:"pin" example:"123456"`
	Tel   string `json:"tel" example:"0812345678"`
	RefNo string `json:"refNo" example:"ABC123"`
}

type VerifyOTPResponse struct {
	Status string `json:"status" example:"success"`
}

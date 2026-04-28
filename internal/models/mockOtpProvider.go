package models

import "time"

type OTPStatus string

const (
	StatusPending  OTPStatus = "pending"
	StatusVerified OTPStatus = "verified"
	StatusExpired  OTPStatus = "expired"
)

type MockOTP struct {
	ID uint `gorm:"primaryKey"`

	Token string `gorm:"size:100;uniqueIndex;not null"`
	RefNo string `gorm:"size:50;index;not null"`

	Tel string `gorm:"size:15;index;not null"`

	Pin string `gorm:"size:10;not null"`

	Status OTPStatus `gorm:"size:20;default:'pending'"`

	ExpiresAt time.Time `gorm:"not null"`

	CreatedAt time.Time
}

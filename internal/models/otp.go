package models

import "time"

type Otp struct {
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

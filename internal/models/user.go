package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `gorm:"size:100;not null" json:"name"`
	Email    string `gorm:"uniqueIndex;not null" json:"email"`
	GoogleID string `gorm:"uniqueIndex" json:"google_id"`

	IsValidated bool   `gorm:"default:false" json:"is_validated"`
	Tel         string `gorm:"size:20" json:"tel"`

	CompanyID *uint
	Company   *Company `gorm:"foreignKey:CompanyID"`
}

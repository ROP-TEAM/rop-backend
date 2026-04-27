package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string `gorm:"size:100;not null" json:"name"`
	Email       string `gorm:"uniqueIndex;not null" json:"email"`
	Pictures    string `json:"pictures"`
	GoogleID    string `gorm:"uniqueIndex" json:"google_id"`
	Tel         string `gorm:"type:varchar(10);uniqueIndex" json:"tel"`
	IsValidated bool   `gorm:"default:false" json:"is_validated"`
}

package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name        string  `gorm:"size:100;not null" json:"name"`
	Email       string  `gorm:"uniqueIndex;not null" json:"email"`
	GoogleID    string  `gorm:"uniqueIndex" json:"google_id"`
	IsValidated bool    `gorm:"default:false" json:"is_validated"`
	Tel         string  `gorm:"type:varchar(10)" json:"tel"`
	TelOTP      string  `gorm:"size:6" json:"tel_otp"`
	CompanyID   string  `gorm:"column:user_company_fk;index;type:uuid"`
	Company     Company `gorm:"foreignKey:CompanyID"`
	Pictures    string  `json:"picture"`
}

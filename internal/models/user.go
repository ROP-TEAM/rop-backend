package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name             string  `gorm:"type:varchar(100);not null" json:"name"`
	Email            string  `gorm:"uniqueIndex;not null" json:"email"`
	GoogleID         string  `gorm:"uniqueIndex" json:"google_id"`
	IsValidated      bool    `gorm:"default:false" json:"is_validated"`
	IsNeedOnBoarding bool    `gorm:"default:true" json:"is_need_on_boarding"`
	Tel              string  `gorm:"type:varchar(10)" json:"tel"`
	CompanyID        *string `gorm:"column:user_company_fk;index;type:uuid"`
	Company          Company `gorm:"foreignKey:CompanyID"`
}

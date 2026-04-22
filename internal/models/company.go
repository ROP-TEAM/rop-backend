package models

import "gorm.io/gorm"

type Company struct {
	gorm.Model
	ID          uint    `gorm:"primaryKey" json:"id"`
	Name        string  `gorm:"size:100;not null" json:"name"`
	Type        string  `gorm:"size:50;not null" json:"type"`
	Province    string  `gorm:"size:50;not null" json:"province"`
	District    string  `gorm:"size:50;not null" json:"district"`
	SubDistrict string  `gorm:"size:50;not null" json:"sub_district"`
	Address     string  `gorm:"size:255;not null" json:"address"`
	Alley       *string `gorm:"size:255" json:"alley"`
	PostalCode  string  `gorm:"size:20;not null" json:"postal_code"`

	Tel string `gorm:"size:20" json:"tel"`
}

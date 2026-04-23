package models

type Company struct {
	ID   uint   `gorm:"primaryKey" json:"id"`
	Name string `gorm:"type:varchar(100);" json:"name"`
	Type string `gorm:"size:50;not null" json:"type"`

	Address     string  `gorm:"type:text" json:"address"`
	District    string  `gorm:"type:varchar(100)" json:"district"`
	SubDistrict string  `gorm:"type:varchar(100)" json:"sub_district"`
	Province    string  `gorm:"type:varchar(100)" json:"province"`
	PostalCode  string  `gorm:"type:varchar(10)" json:"postal_code"`
	Alley       *string `gorm:"type:varchar(100)" json:"alley"`

	Tel string `gorm:"type:varchar(20)" json:"tel"`
}

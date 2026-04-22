package models

type Company struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"type:varchar(100);" json:"name"`

	Address     string  `gorm:"type:text" json:"address"`
	District    string  `gorm:"type:varchar(100)" json:"district"`
	SubDistrict string  `gorm:"type:varchar(100)" json:"sub_district"`
	Province    string  `gorm:"type:varchar(100)" json:"province"`
	PostalCode  string  `gorm:"type:varchar(10)" json:"postal_code"`
	Alley       *string `gorm:"type:varchar(100)" json:"alley"`

	Tel string `gorm:"type:varchar(20)" json:"tel"`
}

package models

type Company struct {
	ID   string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()" json:"id"`
	Name string `gorm:"type:varchar(100);" json:"name"`

	Address     string `gorm:"type:text" json:"address"`
	District    string `gorm:"type:varchar(100)" json:"district"`
	SubDistrict string `gorm:"type:varchar(100)" json:"sub_district"`
	Province    string `gorm:"type:varchar(100)" json:"province"`
	PostalCode  string `gorm:"type:char(5)" json:"postal_code"`
	Type        string `gorm:"type:varchar(100);not null" json:"type"`
	Tel         string `gorm:"type:varchar(10)" json:"tel"`
}

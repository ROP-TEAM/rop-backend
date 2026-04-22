package models

type Client struct {
	ID        uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string `gorm:"size:255;not null" json:"name"`
	Tel       string `gorm:"size:10" json:"tel"`
	CompanyID string `gorm:"column:client_company_fk;not null;index;type:uuid" json:"company_id"`
}

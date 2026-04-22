package models

type DistributionPoint struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	Latitude  float64 `gorm:"type:decimal(9,6)" json:"latitude"`
	Longitude float64 `gorm:"type:decimal(9,6)" json:"longitude"`
	CompanyID string  `gorm:"column:distribution_point_company_fk;not null;index;type:uuid" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID;references:ID" json:"-"`
}

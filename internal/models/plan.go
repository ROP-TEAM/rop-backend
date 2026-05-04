package models

import (
	"time"

	"gorm.io/gorm"
)

type Plan struct {
	ID        string         `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name      string         `gorm:"type:varchar(100)" json:"name"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
	PlanDate  time.Time      `gorm:"type:date;index" json:"plane_date"`
	Status    string         `gorm:"type:varchar(20);default:'pending';check:status IN ('pending', 'processing', 'completed', 'failed')"`
	CompanyID string         `gorm:"column:plan_company_fk;type:uuid;index;" json:"company_id"`
	Company   Company        `gorm:"foreignKey:CompanyID"`
	Routes    []Route        `json:"routes"`
}

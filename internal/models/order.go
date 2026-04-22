package models

import (
	"time"

	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Name   string `gorm:"size:255;not null" json:"name"`
	Note   string `gorm:"type:text" json:"note"`
	Status string `gorm:"index;type:varchar(20);not null;check:status IN ('done', 'pending', 'cancel');default:'pending'" json:"status"`
	// 1:for pickeup, 0 for delivery
	Mode int `gorm:"type:smallint;not null;default:0;check:mode IN (0,1)" json:"mode"`

	WeightKg    float64 `gorm:"type:decimal(10,2)" json:"weight_kg"`
	VolumeM3    float64 `gorm:"type:decimal(10,2)" json:"volume_m3"`
	ServiceTime int     `gorm:"default:5" json:"service_time"`
	// 3: Critical, 2: High, 1: Medium, 0: Low
	Priority int        `gorm:"type:smallint;not null;default:1;check:priority BETWEEN 0 AND 3" json:"priority"`
	Skills   []TagSkill `gorm:"many2many:order_tag_skills;foreignKey:ID;joinForeignKey:OrderID;References:ID;joinReferences:TagSkillID"`

	TimeWindowStart time.Time `gorm:"type:time" json:"time_window_start"`
	TimeWindowEnd   time.Time `gorm:"type:time" json:"time_window_end"`
	DateWindowStart time.Time `gorm:"index;type:date" json:"date_window_start"`
	DateWindowEnd   time.Time `gorm:"type:date" json:"date_window_end"`

	DesLatitude  float64 `gorm:"type:decimal(9,6);not null" json:"des_latitude"`
	DesLongitude float64 `gorm:"type:decimal(9,6);not null" json:"des_longitude"`

	CompanyID string  `gorm:"column:order_company_fk;not null;index;type:uuid" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID" json:"-"`
	UserID    uint    `gorm:"column:order_user_fk;not null;index" json:"user_id"`
	User      User    `gorm:"foreignKey:UserID" json:"-"`
	ClientID  uint    `gorm:"column:order_client_fk;not null" json:"client_id"`
	Client    Client  `gorm:"foreignKey:ClientID" json:"client"`
}

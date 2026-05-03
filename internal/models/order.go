package models

import (
	"gorm.io/gorm"
)

type Order struct {
	gorm.Model
	Name string `gorm:"type:varchar(255);not null" json:"name"`
	Note string `gorm:"type:text" json:"note"`
	// 1:for pickeup, 0 for delivery
	Type int `gorm:"type:smallint;not null;default:0;check:type IN (0,1)" json:"type"`

	Capacity    float64 `gorm:"type:numeric(10,2);not null" json:"capacity"`
	ServiceTime *int    `gorm:"default:10" json:"service_time"`
	// 3: Critical, 2: High, 1: Medium, 0: Low
	Priority int        `gorm:"type:smallint;not null;default:0;check:priority BETWEEN 0 AND 3" json:"priority"`
	Skills   []TagSkill `gorm:"many2many:order_tag_skills;foreignKey:ID;joinForeignKey:OrderID;References:ID;joinReferences:TagSkillID"`

	TimeWindowStart *int `gorm:"" json:"time_window_start"`
	TimeWindowEnd   *int `gorm:"" json:"time_window_end"`

	DesLatitude  *float64 `gorm:"type:double precision;not null" json:"des_latitude"`
	DesLongitude *float64 `gorm:"type:double precision;not null" json:"des_longitude"`

	PlanID string `gorm:"column:order_plan_fk;type:uuid;index" json:"plan_id"`
	Plan   Plan   `gorm:"foreignKey:PlanID"`
}

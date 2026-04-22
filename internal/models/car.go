package models

import "time"

type Car struct {
	ID uint `gorm:"primaryKey"`

	DisplayName string `gorm:"size:100" json:"display_name"`
	PlateNumber string `gorm:"size:20;not null;uniqueIndex:idx_company_plate" json:"plate_number"`
	Model       string `gorm:"size:100" json:"model"`

	MaxWeightKg float64 `gorm:"type:decimal(10,2);not null" json:"max_weight_kg"`
	MaxVolumeM3 float64 `gorm:"type:decimal(10,2);not null" json:"max_volume_m3"`
	MaxTask     int     `gorm:"default:20" json:"max_task"`

	DailyBreakTimeStart time.Time `gorm:"type:time" json:"daily_break_time_start"`
	DailyBreakTimeEnd   time.Time `gorm:"type:time" json:"daily_break_time_end"`
	DailyWorkTimeStart  time.Time `gorm:"type:time" json:"daily_work_time_start"`
	DailyWorkTimeEnd    time.Time `gorm:"type:time" json:"daily_work_time_end"`

	CompanyID string  `gorm:"column:company_car_fk;not null;index;type:uuid;uniqueIndex:idx_company_plate" json:"company_id"`
	Company   Company `gorm:"foreignKey:CompanyID;references:ID" json:"company"`

	DefaultDepotID uint              `gorm:"column:default_depot_fk" json:"default_depot_id"`
	DefaultDepot   DistributionPoint `gorm:"foreignKey:DefaultDepotID" json:"default_depot"`
	CurrentLat     float64           `gorm:"type:decimal(9,6)" json:"current_latitude"`
	CurrentLon     float64           `gorm:"type:decimal(9,6)" json:"current_lon"`

	// many-to-many relatioinship
	Skills []TagSkill `gorm:"many2many:car_tag_skills;foreignKey:ID;joinForeignKey:CarID;References:ID;joinReferences:TagSkillID" json:"skills"`

	//car is ready to be used or not
	IsReady bool `gorm:"default:true" json:"is_ready"`
}

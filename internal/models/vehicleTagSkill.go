package models

type VehicleTagSkill struct {
	VehicleID  uint `gorm:"primaryKey" json:"vehicle_id"`
	TagSkillID uint `gorm:"primaryKey" json:"tag_skill_id"`
}

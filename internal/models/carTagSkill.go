package models

type CarTagSkill struct {
	CarID      uint `gorm:"primaryKey" json:"car_id"`
	TagSkillID uint `gorm:"primaryKey" json:"tag_skill_id"`
}

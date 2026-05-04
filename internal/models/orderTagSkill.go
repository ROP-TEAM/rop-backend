package models

type OrderTagSkill struct {
	OrderID    uint `gorm:"primaryKey" json:"order_id"`
	TagSkillID uint `gorm:"primaryKey" json:"tag_skill_id"`
}

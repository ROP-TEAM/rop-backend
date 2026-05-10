package models

type TagSkill struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"type:varchar(50);not null;uniqueIndex:idx_name_plan" json:"name"`
	Color  string `gorm:"type:varchar(20);" json:"color"`
	PlanID string `gorm:"column:tag_skill_plan_fk;type:uuid;uniqueIndex:idx_name_plan" json:"plan_id"`
	Plan   Plan   `gorm:"foreignKey:PlanID"`
}

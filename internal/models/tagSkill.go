package models

type TagSkill struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name   string `gorm:"type:varchar(50);not null;unique" json:"name"`
	Color  string `gorm:"type:varchar(20);" json:"color"`
	PlanID string `gorm:"column:tag_skill_plan_fk;type:uuid;index" json:"plan_id"`
	Plan   Plan   `gorm:"foreignKey:PlanID"`
}

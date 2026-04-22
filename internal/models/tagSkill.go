package models

type TagSkill struct {
	ID   uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name string `gorm:"size:50;not null;unique" json:"name"`
}

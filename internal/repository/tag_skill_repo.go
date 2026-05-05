package repository

import (
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

type TagSkillRepository struct {
	db *gorm.DB
}

func NewTagSkillRepository(db *gorm.DB) *TagSkillRepository {
	return &TagSkillRepository{db: db}
}

func (r *TagSkillRepository) Create(skills []models.TagSkill) error {
	return r.db.Create(&skills).Error
}

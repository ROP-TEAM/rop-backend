package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type OrderTagSkillRepository struct {
	db *gorm.DB
}

func NewOrderTagSkillRepository(db *gorm.DB) *OrderTagSkillRepository {
	return &OrderTagSkillRepository{db: db}
}

func (r *OrderTagSkillRepository) BatchCreate(ctx context.Context, records []models.OrderTagSkill) error {
	return r.db.WithContext(ctx).CreateInBatches(records, 100).Error
}

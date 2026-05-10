package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type VehicleTagSkillRepository struct {
	db *gorm.DB
}

func (r *VehicleTagSkillRepository) WithTx(tx *gorm.DB) *VehicleTagSkillRepository {
	return &VehicleTagSkillRepository{db: tx}
}

func NewVehicleTagSkillRepository(db *gorm.DB) *VehicleTagSkillRepository {
	return &VehicleTagSkillRepository{db: db}
}

func (r *VehicleTagSkillRepository) BatchCreate(ctx context.Context, records []models.VehicleTagSkill) error {
	return r.db.WithContext(ctx).CreateInBatches(records, 100).Error
}

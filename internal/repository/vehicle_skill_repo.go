package repository

import (
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

type VehicleSkillRepository struct {
	db *gorm.DB
}

func NewVehicleSkillRepository(db *gorm.DB) *VehicleSkillRepository {
	return &VehicleSkillRepository{db: db}
}

func (r *VehicleSkillRepository) Replace(
	vehicleID uint,
	skills []models.VehicleTagSkill,
) error {
	tx := r.db.Begin()

	if err := tx.Where("Vehicle_id = ?", vehicleID).Delete(&models.VehicleTagSkill{}).
		Error; err != nil {
		tx.Rollback()
	}

	if len(skills) > 0 {
		if err := tx.Create(&skills).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

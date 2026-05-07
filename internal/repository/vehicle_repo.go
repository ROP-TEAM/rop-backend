package repository

import (
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(
	vehicle *models.Vehicle,
) error {
	return r.db.Create(vehicle).Error
}

func (r *VehicleRepository) Update(
	vehicle *models.Vehicle,
) error {
	return r.db.Save(vehicle).Error
}

func (r *VehicleRepository) FindByIDAndPlan(
	id uint,
	planID string,
) (*models.Vehicle, error) {

	var vehicle models.Vehicle

	err := r.db.
		Where("id = ? AND vehicle_plan_fk = ?", id, planID).
		First(&vehicle).Error

	if err != nil {
		return nil, err
	}

	return &vehicle, nil
}

func (r *VehicleRepository) CreateSkills(
	skills []models.VehicleTagSkill,
) error {
	return r.db.Create(&skills).Error
}

func (r *VehicleRepository) CountByPlan(
	planID string,
	count *int64,
) error {
	return r.db.
		Model(&models.Vehicle{}).
		Where("vehicle_plan_fk = ?", planID).
		Count(count).
		Error
}

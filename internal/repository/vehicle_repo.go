package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func (r *VehicleRepository) WithTx(tx *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: tx}
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(
	vehicle *models.Vehicle,
) error {
	return r.db.Create(vehicle).Error
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

	if err := r.db.
		Where(
			"id = ? AND vehicle_plan_fk = ?",
			id,
			planID,
		).
		First(&vehicle).Error; err != nil {

		return nil, err
	}

	return &vehicle, nil
}

func (r *VehicleRepository) DeleteSkills(
	vehicleID uint,
) error {

	return r.db.
		Where("vehicle_id = ?", vehicleID).
		Delete(&models.VehicleTagSkill{}).
		Error
}

func (r *VehicleRepository) FindByIDWithSkills(
	id uint,
) (*models.Vehicle, error) {

	var vehicle models.Vehicle

	if err := r.db.
		Preload("Skills").
		Where("id = ?", id).
		First(&vehicle).Error; err != nil {

		return nil, err
	}

	return &vehicle, nil
}

// กันส่งชื่อที่มีแล้วอีกทีเพื่อความชัวร์
func (r *VehicleRepository) FindTagSkillByName(
	name string,
) (*models.TagSkill, error) {

	var tag models.TagSkill

	err := r.db.
		Where(
			"name = ?",
			name,
		).
		First(&tag).Error

	if err != nil {
		return nil, err
	}

	return &tag, nil
}

func (r *VehicleRepository) CreateTagSkill(
	tag *models.TagSkill,
) error {
	return r.db.Create(tag).Error
}

func (r *VehicleRepository) Delete(
	id uint,
) error {

	return r.db.
		Delete(&models.Vehicle{}, id).
		Error
}

func (r *VehicleRepository) BatchCreate(ctx context.Context, vehicles []models.Vehicle) error {
	return r.db.WithContext(ctx).CreateInBatches(vehicles, 100).Error
}

func (r *VehicleRepository) CountOrders(
	tagSkillID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&models.OrderTagSkill{}).
		Where(
			"tag_skill_id = ?",
			tagSkillID,
		).
		Count(&count).
		Error

	return count, err
}

func (r *VehicleRepository) CountVehicles(
	tagSkillID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Model(&models.VehicleTagSkill{}).
		Where(
			"tag_skill_id = ?",
			tagSkillID,
		).
		Count(&count).
		Error

	return count, err
}

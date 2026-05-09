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

func (r *TagSkillRepository) Create(
	skill *models.TagSkill,
) error {

	return r.db.Create(skill).Error
}

func (r *TagSkillRepository) FindByIDAndPlan(
	id uint,
	planID string,
) (*models.TagSkill, error) {

	var skill models.TagSkill

	err := r.db.
		First(
			&skill,
			"id = ? AND tag_skill_plan_fk = ?",
			id,
			planID,
		).Error

	if err != nil {
		return nil, err
	}

	return &skill, nil
}

func (r *TagSkillRepository) Update(
	skill *models.TagSkill,
) error {

	return r.db.Save(skill).Error
}

func (r *TagSkillRepository) CountOrders(
	skillID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Table("order_tag_skills").
		Where("tag_skill_id = ?", skillID).
		Count(&count).Error

	return count, err
}

func (r *TagSkillRepository) CountVehicles(
	skillID uint,
) (int64, error) {

	var count int64

	err := r.db.
		Table("vehicle_tag_skills").
		Where("tag_skill_id = ?", skillID).
		Count(&count).Error

	return count, err
}

func (r *TagSkillRepository) CheckPlanExists(planID string) error {
	var plan models.Plan

	return r.db.First(&plan, "id = ?", planID).Error
}

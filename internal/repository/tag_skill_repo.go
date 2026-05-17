package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type TagSkillRepository struct {
	db *gorm.DB
}

func (r *TagSkillRepository) WithTx(tx *gorm.DB) *TagSkillRepository {
	return &TagSkillRepository{db: tx}
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

func (r *TagSkillRepository) CreateWithContext(ctx context.Context, skill *models.TagSkill) error {
	return r.db.WithContext(ctx).Create(skill).Error
}

func (r *TagSkillRepository) BatchCreate(ctx context.Context, skills []models.TagSkill) error {
	return r.db.WithContext(ctx).Create(&skills).Error
}

func (r *TagSkillRepository) FindByName(
	name string,
) (*models.TagSkill, error) {

	var skill models.TagSkill

	err := r.db.
		Where("name = ?", name).
		First(&skill).
		Error

	if err != nil {
		return nil, err
	}

	return &skill, nil
}

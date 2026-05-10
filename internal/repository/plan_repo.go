package repository

import (
	"ROP_Backend/internal/models"
	"context"
	"errors"

	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func (r *PlanRepository) WithTx(tx *gorm.DB) *PlanRepository {
	return &PlanRepository{db: tx}
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(ctx context.Context, req *models.Plan) error {
	return r.db.WithContext(ctx).Create(&req).Error
}

func (r *PlanRepository) UpdateByID(ctx context.Context, planID string, companyID string, req *models.Plan) error {
	result := r.db.WithContext(ctx).
		Model(req).
		Where("id = ? AND plan_company_fk = ?", planID, companyID).
		Updates(req)

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

func (r *PlanRepository) HardDeleteByID(ctx context.Context, planID string, companyID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var plan models.Plan
		err := tx.Unscoped().Where("id = ? AND plan_company_fk = ?", planID, companyID).
			First(&plan).Error
		if err != nil {
			return err
		}

		err = tx.Exec("DELETE FROM order_tag_skills WHERE order_id IN (SELECT id FROM orders WHERE order_plan_fk = ?)", planID).Error
		if err != nil {
			return err
		}

		err = tx.Exec("DELETE FROM vehicle_tag_skills WHERE vehicle_id IN (SELECT id FROM vehicles WHERE vehicle_plan_fk = ?)", planID).Error
		if err != nil {
			return err
		}

		err = tx.Exec("DELETE FROM stops WHERE route_id IN (SELECT id FROM routes WHERE plan_id = ?)", planID).Error
		if err != nil {
			return err
		}

		// first child
		err = tx.Unscoped().Where("plan_id = ?", planID).Delete(&models.Route{}).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Where("vehicle_plan_fk = ?", planID).Delete(&models.Vehicle{}).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Where("order_plan_fk = ?", planID).Delete(&models.Order{}).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Where("tag_skill_plan_fk = ?", planID).Delete(&models.TagSkill{}).Error
		if err != nil {
			return err
		}

		err = tx.Unscoped().Delete(&plan).Error
		if err != nil {
			return err
		}

		return nil
	})
}

func (r *PlanRepository) FindByID(ctx context.Context, planID string, companyID string) (*models.Plan, error) {
	var plan models.Plan
	err := r.db.WithContext(ctx).
		Preload("TagSkills").
		Preload("Orders.Skills").
		Preload("Vehicles.Skills").
		Preload("Routes.Stops").
		Where("id = ? AND plan_company_fk = ?", planID, companyID).
		First(&plan).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrPlanNotFound
	}
	return &plan, err
}

func (r *PlanRepository) CountExistingNameCopies(ctx context.Context, companyID string, originalName string) (*int64, error) {
	var count int64
	pattern := originalName + " (Copy%"

	err := r.db.WithContext(ctx).Model(&models.Plan{}).
		Where("plan_company_fk = ? AND name LIKE ?", companyID, pattern).
		Count(&count).Error

	return &count, err
}

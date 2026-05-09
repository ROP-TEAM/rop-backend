package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
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

func (r *PlanRepository) DeleteByID(ctx context.Context, planID string, companyID string) error {
	result := r.db.WithContext(ctx).
		Where("id = ? AND plan_company_fk = ?", planID, companyID).
		Delete(&models.Plan{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return ErrNoRowsAffected
	}

	return nil
}

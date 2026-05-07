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

func (r *PlanRepository) Create(ctx context.Context, req *models.Plan) (*models.Plan, error) {
	err := r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return nil, err
	}
	return req, err
}

func (r *PlanRepository) UpdateByID(ctx context.Context, id string, req *models.Plan) (*models.Plan, error) {
	var plan models.Plan
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&plan).Error; err != nil {
		return nil, err
	}

	result := r.db.WithContext(ctx).
		Model(&plan).
		Updates(req)

	if result.Error != nil {
		return nil, result.Error
	}

	if result.RowsAffected == 0 {
		return nil, ErrNoRowsAffected
	}

	return &plan, nil
}

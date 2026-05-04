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
	return &PlanRepository{}
}

func (r *PlanRepository) Create(ctx context.Context, req *models.Plan) (*models.Plan, error) {
	err := r.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return nil, err
	}
	return req, err
}

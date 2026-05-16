package repository

import (
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

type PlanRepository struct {
	db *gorm.DB
}

func NewPlanRepository(db *gorm.DB) *PlanRepository {
	return &PlanRepository{db: db}
}

func (r *PlanRepository) Create(plan *models.Plan) error {
	return r.db.Create(plan).Error
}

func (r *PlanRepository) FindByID(id string) (*models.Plan, error) {
	var plan models.Plan
	if err := r.db.First(&plan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) FindByIDWithRoutes(id string) (*models.Plan, error) {
	var plan models.Plan
	if err := r.db.
		Preload("Routes.Stops").
		First(&plan, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &plan, nil
}

func (r *PlanRepository) FindByCompanyID(companyID string) ([]models.Plan, error) {
	var plans []models.Plan
	if err := r.db.
		Where("plan_company_fk = ?", companyID).
		Order("created_at DESC").
		Find(&plans).Error; err != nil {
		return nil, err
	}
	return plans, nil
}

func (r *PlanRepository) UpdateStatus(id string, status string) error {
	return r.db.
		Model(&models.Plan{}).
		Where("id = ?", id).
		Update("status", status).Error
}

func (r *PlanRepository) SaveRoutes(routes []models.Route) error {
	return r.db.Create(&routes).Error
}

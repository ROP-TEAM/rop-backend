package repository

import (
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(orders *models.Order) error {
	return r.db.Create(&orders).Error
}

func (r *OrderRepository) CreateSkills(
	skills []models.OrderTagSkill,
) error {
	return r.db.Create(&skills).Error
}

func (r *OrderRepository) CountByPlan(
	planID string,
	count *int64,
) error {
	return r.db.Model(&models.Order{}).Where("order_plan_fk = ?", planID).
		Count(count).Error
}

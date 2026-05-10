package repository

import (
	"ROP_Backend/internal/models"
	"context"

	"gorm.io/gorm"
)

type OrderRepository struct {
	db *gorm.DB
}

func (r *OrderRepository) WithTx(tx *gorm.DB) *OrderRepository {
	return &OrderRepository{db: tx}
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

func (r *OrderRepository) FindByIDAndPlan(
	id uint,
	planID string,
) (*models.Order, error) {

	var order models.Order

	if err := r.db.
		Where(
			"id = ? AND order_plan_fk = ?",
			id,
			planID,
		).
		First(&order).Error; err != nil {

		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) DeleteSkills(
	orderID uint,
) error {

	return r.db.
		Where("order_id = ?", orderID).
		Delete(&models.OrderTagSkill{}).
		Error
}

func (r *OrderRepository) FindByIDWithSkills(
	id uint,
) (*models.Order, error) {

	var order models.Order

	if err := r.db.
		Preload("Skills").
		Where("id = ?", id).
		First(&order).Error; err != nil {

		return nil, err
	}

	return &order, nil
}

func (r *OrderRepository) Delete(
	id uint,
) error {

	return r.db.
		Delete(&models.Order{}, id).
		Error
}

func (r *OrderRepository) BatchCreate(ctx context.Context, orders []models.Order) error {
	return r.db.WithContext(ctx).CreateInBatches(orders, 100).Error
}

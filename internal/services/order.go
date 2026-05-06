package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/validators"
	"errors"

	"gorm.io/gorm"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(db *gorm.DB) *OrderService {
	repo := repository.NewOrderRepository(db)
	return &OrderService{repo: repo}
}

func (s *OrderService) GroupCreate(req dto.GroupCreateOrder) error {
	if len(req.Orders) == 0 {
		return errors.New("orders is empty")
	}

	var orders []models.Order

	for _, o := range req.Orders {

		order := models.Order{
			Name: o.Name,
			Note: o.Note,
			Type: o.Type,

			Capacity:    o.Capacity,
			ServiceTime: o.ServiceTime,
			Priority:    o.Priority,

			TimeWindowStart: o.TimeWindowStart,
			TimeWindowEnd:   o.TimeWindowEnd,

			DesLatitude:  &o.DesLatitude,
			DesLongitude: &o.DesLongitude,

			PlanID: req.PlanID,
		}

		if err := validators.ValidateOrder(&order); err != nil {
			return err
		}

		orders = append(orders, order)
	}
	return s.repo.Create(orders)
}

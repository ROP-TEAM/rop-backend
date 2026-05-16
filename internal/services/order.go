package services

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/validators"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type OrderService struct {
	repo *repository.OrderRepository
}

func NewOrderService(db *gorm.DB) *OrderService {
	repo := repository.NewOrderRepository(db)
	return &OrderService{repo: repo}
}

func (s *OrderService) GroupCreate(
	req dto.GroupCreateOrder) ([]response.OrderResponse, error) {
	var responses []response.OrderResponse

	var count int64
	const MaxOrders = 200

	if len(req.Orders) == 0 {
		return nil, errors.New("orders is empty")
	}

	if err := s.repo.CountByPlan(
		req.PlanID,
		&count,
	); err != nil {
		return nil, err
	}

	if count+int64(len(req.Orders)) > MaxOrders {
		return nil, fmt.Errorf(
			"maximum %d orders per plan",
			MaxOrders,
		)
	}

	for _, o := range req.Orders {

		newOrder := models.Order{
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

		if err := validators.ValidateOrder(&newOrder); err != nil {
			return nil, err
		}

		if err := s.repo.Create(&newOrder); err != nil {
			return nil, err
		}

		var skills []models.OrderTagSkill

		for _, skillID := range o.TagSkillID {
			skills = append(skills, models.OrderTagSkill{
				OrderID:    newOrder.ID,
				TagSkillID: skillID,
			})
		}

		if len(skills) > 0 {
			if err := s.repo.CreateSkills(skills); err != nil {
				return nil, err
			}
		}

		responses = append(responses, response.OrderResponse{
			Name: newOrder.Name,
			Note: newOrder.Note,

			Type: newOrder.Type,

			Capacity: newOrder.Capacity,

			ServiceTime: newOrder.ServiceTime,
			Priority:    newOrder.Priority,

			TagSkillID: o.TagSkillID,
		})

	}
	return responses, nil
}

func (s *OrderService) Update(
	orderID uint,
	req dto.UpdateOrder,
) (*response.OrderResponse, error) {

	order, err := s.repo.FindByIDAndPlan(
		orderID,
		req.PlanID,
	)

	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		order.Name = *req.Name
	}

	if req.Note != nil {
		order.Note = *req.Note
	}

	if req.Type != nil {
		order.Type = *req.Type
	}

	if req.Capacity != nil {
		order.Capacity = *req.Capacity
	}

	if req.ServiceTime != nil {
		order.ServiceTime = req.ServiceTime
	}

	if req.Priority != nil {
		order.Priority = *req.Priority
	}

	if req.TimeWindowStart != nil {
		order.TimeWindowStart = req.TimeWindowStart
	}

	if req.TimeWindowEnd != nil {
		order.TimeWindowEnd = req.TimeWindowEnd
	}

	if req.DesLatitude != nil {
		order.DesLatitude = req.DesLatitude
	}

	if req.DesLongitude != nil {
		order.DesLongitude = req.DesLongitude
	}

	if err := validators.ValidateOrder(order); err != nil {
		return nil, err
	}

	if err := s.repo.Update(order); err != nil {
		return nil, err
	}

	if req.TagSkillID != nil {

		if err := s.repo.DeleteSkills(order.ID); err != nil {
			return nil, err
		}

		var skills []models.OrderTagSkill

		for _, skillID := range *req.TagSkillID {

			skills = append(skills, models.OrderTagSkill{
				OrderID:    order.ID,
				TagSkillID: skillID,
			})
		}

		if len(skills) > 0 {
			if err := s.repo.CreateSkills(skills); err != nil {
				return nil, err
			}
		}
	}

	updatedOrder, err := s.repo.FindByIDWithSkills(
		order.ID,
	)

	if err != nil {
		return nil, err
	}

	var tagSkillIDs []uint

	for _, s := range updatedOrder.Skills {
		tagSkillIDs = append(
			tagSkillIDs,
			s.ID,
		)
	}

	responseData := response.OrderResponse{
		Name: updatedOrder.Name,
		Note: updatedOrder.Note,

		Type: updatedOrder.Type,

		Capacity: updatedOrder.Capacity,

		ServiceTime: updatedOrder.ServiceTime,
		Priority:    updatedOrder.Priority,

		TagSkillID: tagSkillIDs,
	}

	return &responseData, nil
}

func (s *OrderService) Delete(req dto.DeleteOrder) error {
	for _, id := range req.ID {
		order, err := s.repo.FindByIDAndPlan(
			id,
			req.PlanID,
		)

		if err != nil {
			return err
		}

		if err := s.repo.DeleteSkills(
			order.ID,
		); err != nil {
			return err
		}

		if err := s.repo.Delete(
			order.ID,
		); err != nil {
			return err
		}
	}

	return nil
}

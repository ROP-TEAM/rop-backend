package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/validators"
)

type PlanService struct {
}

func NewPlanService() *PlanService {
	return &PlanService{}
}

func (s *PlanService) CreatePlan(req dto.PlanRequest) error {
	if err := validators.ValidatePlan(req); err != nil {
		return err
	}

	// call algorithm here in the future

	return nil
}

package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"context"
	"errors"
	"log"

	"gorm.io/gorm"
)

var (
	ErrInvalidUser = errors.New("invalid user or user not found ")
)

type PlanService struct {
	planRepository *repository.PlanRepository
	userRepository *repository.UserRepository
}

func NewPlanService(db *gorm.DB) *PlanService {
	planRepo := repository.NewPlanRepository(db)
	userRepo := repository.NewUserRepository(db)
	return &PlanService{
		planRepository: planRepo,
		userRepository: userRepo,
	}
}

func (s *PlanService) CreateByUserID(ctx context.Context, userID uint, req *dto.CreatePlanRequest) (*dto.CreatePlanResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		log.Printf("[planService]: finding company id by user id: %v", err)
		return nil, ErrInvalidUser
	}

	plan, err := s.planRepository.Create(ctx, &models.Plan{
		CompanyID: *user.CompanyID,
		Name:      req.Name,
	})
	if err != nil {
		log.Printf("[planService]: creating plan: %v", err)
		return nil, err
		//internal error ไป
	}

	return &dto.CreatePlanResponse{
		Name: plan.Name,
	}, nil
}

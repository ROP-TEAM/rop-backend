package services

import (
	"ROP_Backend/internal/config"
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
	cfg            *config.Config
}

func NewPlanService(db *gorm.DB, cfg *config.Config) *PlanService {
	planRepo := repository.NewPlanRepository(db)
	userRepo := repository.NewUserRepository(db)
	return &PlanService{
		planRepository: planRepo,
		userRepository: userRepo,
		cfg:            cfg,
	}
}

func (s *PlanService) CreateByUserID(ctx context.Context, userID uint, req *dto.CreatePlanRequest) (*dto.CreatePlanResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		log.Printf("[planService]: finding company id by user id: %v", err)
		return nil, ErrInvalidUser
	}

	newPlan := models.Plan{
		CompanyID: *user.CompanyID,
		Name:      *req.Name,
		PlanDate:  req.PlanDate,
	}
	err = s.planRepository.Create(ctx, &newPlan)
	if err != nil {
		log.Printf("[planService]: creating plan: %v", err)
		return nil, err
	}

	return &dto.CreatePlanResponse{
		ID:        newPlan.ID,
		Name:      newPlan.Name,
		CreatedAt: newPlan.CreatedAt,
	}, nil
}

func (s *PlanService) UpdateNameByID(ctx context.Context, planID string, userID uint, req *dto.UpdatePlanNameByIDRequest) (*dto.UpdatePlanNameByIDResponse, error) {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, err
	}

	if user.CompanyID == nil {
		return nil, errors.New("user has no company")
	}

	updatedPlan := models.Plan{Name: *req.Name}
	err = s.planRepository.UpdateByID(ctx, planID, *user.CompanyID, &updatedPlan)
	if err != nil {
		return nil, err
	}

	return &dto.UpdatePlanNameByIDResponse{Name: updatedPlan.Name, UpdatedAt: updatedPlan.UpdatedAt}, nil
}

func (s *PlanService) DeleteByID(ctx context.Context, planID string, userID uint) error {
	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return err
	}

	if user.CompanyID == nil {
		return errors.New("user has no company")
	}

	return s.planRepository.DeleteByID(ctx, planID, *user.CompanyID)
}

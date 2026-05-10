package services

import (
	"ROP_Backend/internal/config"
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"context"
	"errors"
	"fmt"
	"log"

	"gorm.io/gorm"
)

var (
	ErrInvalidUser = errors.New("invalid user or user not found ")
)

type PlanService struct {
	db                        *gorm.DB //dbนี่ใช้เปิดtransactionอย่างเดียวคคับ หยวนๆกันหน่อย
	planRepository            *repository.PlanRepository
	userRepository            *repository.UserRepository
	tagSkillRepository        *repository.TagSkillRepository
	vehicleRepository         *repository.VehicleRepository
	orderRepository           *repository.OrderRepository
	orderTagSkillRepository   *repository.OrderTagSkillRepository
	vehicleTagSkillRepository *repository.VehicleTagSkillRepository
	cfg                       *config.Config
}

func NewPlanService(db *gorm.DB, cfg *config.Config) *PlanService {
	planRepo := repository.NewPlanRepository(db)
	userRepo := repository.NewUserRepository(db)
	tagSkillRepo := repository.NewTagSkillRepository(db)
	vehicleRepo := repository.NewVehicleRepository(db)
	vehicleTagSkillRepo := repository.NewVehicleTagSkillRepository(db)
	orderRepo := repository.NewOrderRepository(db)
	orderTagSkillRepo := repository.NewOrderTagSkillRepository(db)

	return &PlanService{
		db:                        db,
		planRepository:            planRepo,
		userRepository:            userRepo,
		tagSkillRepository:        tagSkillRepo,
		vehicleRepository:         vehicleRepo,
		orderRepository:           orderRepo,
		orderTagSkillRepository:   orderTagSkillRepo,
		vehicleTagSkillRepository: vehicleTagSkillRepo,
		cfg:                       cfg,
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

	return s.planRepository.HardDeleteByID(ctx, planID, *user.CompanyID)
}

func (s *PlanService) duplicateTagSkills(
	ctx context.Context,
	repo *repository.TagSkillRepository,
	skills []models.TagSkill,
	newPlanID string,
) (map[uint]uint, error) {

	if len(skills) == 0 {
		return make(map[uint]uint), nil
	}

	created := make([]models.TagSkill, len(skills))
	for i, sk := range skills {
		created[i] = models.TagSkill{
			Name:   sk.Name,
			Color:  sk.Color,
			PlanID: newPlanID,
		}
	}

	if err := repo.BatchCreate(ctx, created); err != nil {
		return nil, ErrCreatingTagSkill
	}

	skillMap := make(map[uint]uint)
	for i, sk := range skills {
		skillMap[sk.ID] = created[i].ID
	}

	return skillMap, nil
}

func (s *PlanService) duplicateOrders(
	ctx context.Context,
	orderRepo *repository.OrderRepository,
	orderTagSkillRepo *repository.OrderTagSkillRepository,
	orders []models.Order,
	newPlanID string,
	skillMap map[uint]uint,
) (map[uint]uint, error) {

	newOrders := make([]models.Order, len(orders))
	for i, o := range orders {
		newOrders[i] = models.Order{
			Name:            o.Name,
			Note:            o.Note,
			Type:            o.Type,
			Capacity:        o.Capacity,
			ServiceTime:     o.ServiceTime,
			Priority:        o.Priority,
			TimeWindowStart: o.TimeWindowStart,
			TimeWindowEnd:   o.TimeWindowEnd,
			DesLatitude:     o.DesLatitude,
			DesLongitude:    o.DesLongitude,
			PlanID:          newPlanID,
		}
	}

	if err := orderRepo.BatchCreate(ctx, newOrders); err != nil {
		return nil, ErrCreatingOrder
	}

	orderMap := make(map[uint]uint)
	for i, o := range orders {
		orderMap[o.ID] = newOrders[i].ID
	}

	var orderTagSkills []models.OrderTagSkill
	for i, o := range orders {
		for _, sk := range o.Skills {
			orderTagSkills = append(orderTagSkills, models.OrderTagSkill{
				OrderID:    newOrders[i].ID,
				TagSkillID: skillMap[sk.ID],
			})
		}
	}

	if len(orderTagSkills) > 0 {
		if err := orderTagSkillRepo.BatchCreate(ctx, orderTagSkills); err != nil {
			return nil, ErrCreatingOrderSkill
		}
	}

	return orderMap, nil
}

func (s *PlanService) duplicateVehicles(
	ctx context.Context,
	vehicleRepo *repository.VehicleRepository,
	vehicleTagSkillRepo *repository.VehicleTagSkillRepository,
	vehicles []models.Vehicle,
	newPlanID string,
	skillMap map[uint]uint,
) (map[uint]uint, error) {

	newVehicles := make([]models.Vehicle, len(vehicles))
	for i, v := range vehicles {
		newVehicles[i] = models.Vehicle{
			PlateNumber:         v.PlateNumber,
			Model:               v.Model,
			Name:                v.Name,
			Capacity:            v.Capacity,
			MaxTask:             v.MaxTask,
			DailyBreakTimeStart: v.DailyBreakTimeStart,
			DailyBreakTimeEnd:   v.DailyBreakTimeEnd,
			DailyWorkTimeStart:  v.DailyWorkTimeStart,
			DailyWorkTimeEnd:    v.DailyWorkTimeEnd,
			StartLat:            v.StartLat,
			StartLon:            v.StartLon,
			EndLat:              v.EndLat,
			EndLon:              v.EndLon,
			PlanID:              newPlanID,
		}
	}

	if err := vehicleRepo.BatchCreate(ctx, newVehicles); err != nil {
		return nil, ErrCreatingVehicle
	}

	vehicleMap := make(map[uint]uint)
	for i, o := range vehicles {
		vehicleMap[o.ID] = newVehicles[i].ID
	}

	var vehicleTagSkills []models.VehicleTagSkill
	for i, o := range vehicles {
		for _, sk := range o.Skills {
			vehicleTagSkills = append(vehicleTagSkills, models.VehicleTagSkill{
				VehicleID:  newVehicles[i].ID,
				TagSkillID: skillMap[sk.ID],
			})
		}
	}

	if len(vehicleTagSkills) > 0 {
		if err := vehicleTagSkillRepo.BatchCreate(ctx, vehicleTagSkills); err != nil {
			return nil, ErrCreatingVehicleSkill
		}
	}

	return vehicleMap, nil
}

func (s *PlanService) DuplicateByID(
	ctx context.Context,
	planID string,
	userID uint,
) (*dto.DuplicatePlanByIDResponse, error) {

	user, err := s.userRepository.FindByID(userID)
	if err != nil {
		return nil, ErrUserNotFound
	}

	original, err := s.planRepository.FindByID(ctx, planID, *user.CompanyID)
	if err != nil {
		return nil, ErrPlanNotFound
	}

	if original == nil {
		return nil, ErrPlanNotFound
	}

	var newPlan models.Plan
	err = s.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {

		planRepository := s.planRepository.WithTx(tx)
		skillRepo := s.tagSkillRepository.WithTx(tx)
		orderRepo := s.orderRepository.WithTx(tx)
		orderTagSkillRepo := s.orderTagSkillRepository.WithTx(tx)
		vehicleRepo := s.vehicleRepository.WithTx(tx)
		vehicleTagSkillRepo := s.vehicleTagSkillRepository.WithTx(tx)

		count, err := planRepository.CountExistingNameCopies(ctx, *user.CompanyID, original.Name)
		if err != nil {
			return err
		}

		newPlan = models.Plan{
			Name:      s.formatDuplicateName(original.Name, *count),
			CompanyID: *user.CompanyID,
			PlanDate:  original.PlanDate,
			Status:    "pending",
		}

		err = planRepository.Create(ctx, &newPlan)
		if err != nil {
			return ErrCreatingPlan
		}

		skillMap, err := s.duplicateTagSkills(ctx, skillRepo, original.TagSkills, newPlan.ID)
		if err != nil {
			return err
		}

		_, err = s.duplicateVehicles(ctx, vehicleRepo, vehicleTagSkillRepo, original.Vehicles, newPlan.ID, skillMap)
		if err != nil {
			return err
		}

		_, err = s.duplicateOrders(ctx, orderRepo, orderTagSkillRepo, original.Orders, newPlan.ID, skillMap)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Printf("DuplicatePlan: transaction failed planID=%s: %v", planID, err)
		return nil, err
	}

	return &dto.DuplicatePlanByIDResponse{ID: newPlan.ID, Name: newPlan.Name}, nil
}

func (r *PlanService) formatDuplicateName(baseName string, count int64) string {
	suffix := " (Copy)"
	if count > 0 {
		suffix = fmt.Sprintf(" (Copy %d)", count)
	}
	return baseName + suffix
}

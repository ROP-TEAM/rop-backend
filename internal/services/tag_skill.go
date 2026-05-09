package services

import (
	dto "ROP_Backend/internal/dto/request"
	"ROP_Backend/internal/dto/response"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/validators"
	"errors"

	"gorm.io/gorm"
)

type TagSkillService struct {
	repo *repository.TagSkillRepository
}

func NewTagSkillService(db *gorm.DB) *TagSkillService {
	repo := repository.NewTagSkillRepository(db)
	return &TagSkillService{
		repo: repo,
	}
}

func (s *TagSkillService) GroupCreate(req dto.GroupCreateTagSkill) ([]response.TagSkillResponse, error) {

	if err := s.repo.CheckPlanExists(req.PlanID); err != nil {
		return nil, errors.New("plan not found")
	}

	var responses []response.TagSkillResponse

	for _, sk := range req.Skills {

		if err := validators.ValidateTagSkill(sk); err != nil {
			return nil, err
		}

		var skill models.TagSkill

		if sk.ID == nil {
			newSkill := models.TagSkill{
				Name:   sk.Name,
				Color:  sk.Color,
				PlanID: req.PlanID,
			}

			if err := s.repo.Create(&newSkill); err != nil {
				return nil, err
			}

			skill = newSkill

		} else {

			existing, err := s.repo.FindByIDAndPlan(
				*sk.ID,
				req.PlanID,
			)

			if err != nil {
				return nil, errors.New("skill not found")
			}

			existing.Name = sk.Name
			existing.Color = sk.Color

			if err := s.repo.Update(existing); err != nil {
				return nil, err
			}

			skill = *existing
		}

		orderCount, err := s.repo.CountOrders(skill.ID)
		if err != nil {
			return nil, err
		}

		vehicleCount, err := s.repo.CountVehicles(skill.ID)
		if err != nil {
			return nil, err
		}

		responses = append(responses, response.TagSkillResponse{
			ID:           skill.ID,
			Name:         skill.Name,
			Color:        skill.Color,
			OrderCount:   orderCount,
			VehicleCount: vehicleCount,
		})
	}

	return responses, nil

}

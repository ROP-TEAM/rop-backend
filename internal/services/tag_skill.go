package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/validators"
	"errors"

	"gorm.io/gorm"
)

type TagSkillService struct {
	repo *repository.TagSkillRepository
	db   *gorm.DB
}

func NewTagSkillService(db *gorm.DB) *TagSkillService {
	repo := repository.NewTagSkillRepository(db)
	return &TagSkillService{
		repo: repo,
		db:   db,
	}
}

func (s *TagSkillService) GroupCreate(req dto.GroupCreateTagSkill) error {

	var plan models.Plan
	if err := s.db.First(&plan, "id = ?", req.PlanID).Error; err != nil {
		return errors.New("plan not found")
	}

	var skills []models.TagSkill

	for _, sk := range req.Skills {

		if err := validators.ValidateTagSkill(sk); err != nil {
			return err
		}

		skill := models.TagSkill{
			Name:   sk.Name,
			Color:  sk.Color,
			PlanID: req.PlanID,
		}

		skills = append(skills, skill)
	}

	return s.repo.Create(skills)
}

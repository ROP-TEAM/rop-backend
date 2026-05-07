package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/validators"

	"gorm.io/gorm"
)

type VehicleSkillService struct {
	repo *repository.VehicleSkillRepository
	db   *gorm.DB
}

func NewVehicleSkillService(db *gorm.DB) *VehicleSkillService {

	repo := repository.NewVehicleSkillRepository(db)

	return &VehicleSkillService{
		repo: repo,
		db:   db,
	}
}

func (s *VehicleSkillService) Assign(
	req dto.GroupVehicleSkill,
) error {

	for _, v := range req.Vehicles {

		if err := validators.ValidateVehicleSkill(v); err != nil {
			return err
		}

		var vehicle models.Vehicle

		if err := s.db.
			Where("id = ? AND vehicle_plan_fk = ?", v.VehicleID, req.PlanID).
			First(&vehicle).Error; err != nil {
			return err
		}

		var skills []models.VehicleTagSkill

		for _, skillID := range v.TagSkillID {

			var skill models.TagSkill

			if err := s.db.
				Where("id = ? AND plan_id = ?", skillID, req.PlanID).
				First(&skill).Error; err != nil {
				return err
			}

			skills = append(skills, models.VehicleTagSkill{
				VehicleID:  v.VehicleID,
				TagSkillID: skillID,
			})
		}

		if err := s.repo.Replace(
			v.VehicleID,
			skills,
		); err != nil {
			return err
		}
	}

	return nil
}

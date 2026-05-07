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

type VehicleService struct {
	repo *repository.VehicleRepository
}

func NewVehicleService(db *gorm.DB) *VehicleService {
	repo := repository.NewVehicleRepository(db)
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GroupCreate(
	req dto.GroupCreateVehicle,
) ([]response.VehicleResponse, error) {

	var responses []response.VehicleResponse

	var count int64

	if err := s.repo.CountByPlan(
		req.PlanID,
		&count,
	); err != nil {
		return nil, err
	}

	if count+int64(len(req.Vehicles)) > 30 {
		return nil, errors.New("maximum 30 vehicles per plan")
	}

	for _, v := range req.Vehicles {

		newVehicle := models.Vehicle{
			ProfileID:   v.ProfileID,
			PlateNumber: v.PlateNumber,
			Model:       v.Model,
			Name:        v.Name,
			Capacity:    v.Capacity,
			MaxTask:     v.MaxTask,

			DailyWorkTimeStart:  v.DailyWorkTimeStart,
			DailyWorkTimeEnd:    v.DailyWorkTimeEnd,
			DailyBreakTimeStart: v.DailyBreakTimeStart,
			DailyBreakTimeEnd:   v.DailyBreakTimeEnd,

			StartLat: v.StartLat,
			StartLon: v.StartLon,
			EndLat:   v.EndLat,
			EndLon:   v.EndLon,

			PlanID: req.PlanID,
		}

		if err := validators.ValidateVehicle(&newVehicle); err != nil {
			return nil, err
		}

		if err := s.repo.Create(&newVehicle); err != nil {
			return nil, err
		}

		var skills []models.VehicleTagSkill

		for _, skillID := range v.TagSkillID {

			skills = append(skills, models.VehicleTagSkill{
				VehicleID:  newVehicle.ID,
				TagSkillID: skillID,
			})
		}

		if len(skills) > 0 {
			if err := s.repo.CreateSkills(skills); err != nil {
				return nil, err
			}
		}

		responses = append(responses, response.VehicleResponse{
			VehicleID:   newVehicle.ID,
			ProfileID:   newVehicle.ProfileID,
			PlateNumber: newVehicle.PlateNumber,
			Model:       newVehicle.Model,
			Name:        newVehicle.Name,
			Capacity:    newVehicle.Capacity,
			MaxTask:     newVehicle.MaxTask,
			TagSkillID:  v.TagSkillID,
		})
	}

	return responses, nil
}

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

func (s *VehicleService) mapVehicleResponse(
	vehicle *models.Vehicle,
) (
	response.VehicleResponse,
	error,
) {

	var tagSkills []response.TagSkillResponse

	for _, skill := range vehicle.Skills {

		orderCount, err :=
			s.repo.CountOrders(
				skill.ID,
			)

		if err != nil {
			return response.VehicleResponse{}, err
		}

		vehicleCount, err :=
			s.repo.CountVehicles(
				skill.ID,
			)

		if err != nil {
			return response.VehicleResponse{}, err
		}

		tagSkills = append(
			tagSkills,
			response.TagSkillResponse{
				ID:           skill.ID,
				Name:         skill.Name,
				Color:        skill.Color,
				OrderCount:   orderCount,
				VehicleCount: vehicleCount,
			},
		)
	}

	return response.VehicleResponse{
		VehicleID:   vehicle.ID,
		ProfileID:   vehicle.ProfileID,
		PlateNumber: vehicle.PlateNumber,
		Model:       vehicle.Model,
		Name:        vehicle.Name,
		Capacity:    vehicle.Capacity,
		MaxTask:     vehicle.MaxTask,
		TagSkills:   tagSkills,
	}, nil
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

		for _, tag := range v.TagSkills {

			var tagSkillID uint

			if tag.ID != nil {

				tagSkillID = *tag.ID

			} else {

				existing, err :=
					s.repo.FindTagSkillByName(
						tag.Name,
					)

				if err == nil {

					tagSkillID = existing.ID

				} else {

					newTag := models.TagSkill{
						Name:  tag.Name,
						Color: tag.Color,
					}

					if err := s.repo.CreateTagSkill(
						&newTag,
					); err != nil {
						return nil, err
					}

					tagSkillID = newTag.ID
				}
			}

			skills = append(
				skills,
				models.VehicleTagSkill{
					VehicleID:  newVehicle.ID,
					TagSkillID: tagSkillID,
				},
			)
		}

		if len(skills) > 0 {
			if err := s.repo.CreateSkills(skills); err != nil {
				return nil, err
			}
		}

		createdVehicle, err :=
			s.repo.FindByIDWithSkills(
				newVehicle.ID,
			)

		if err != nil {
			return nil, err
		}

		responseData, err :=
			s.mapVehicleResponse(
				createdVehicle,
			)

		if err != nil {
			return nil, err
		}

		responses = append(
			responses,
			responseData,
		)
	}

	return responses, nil
}

func (s *VehicleService) Update(
	vehicleID uint,
	req dto.UpdateVehicle,
) (*response.VehicleResponse, error) {

	vehicle, err := s.repo.FindByIDAndPlan(
		vehicleID,
		req.PlanID,
	)

	if err != nil {
		return nil, err
	}

	if req.ProfileID != nil {
		vehicle.ProfileID = req.ProfileID
	}

	if req.Name != nil {
		vehicle.Name = *req.Name
	}

	if req.Model != nil {
		vehicle.Model = *req.Model
	}

	if req.PlateNumber != nil {
		vehicle.PlateNumber = *req.PlateNumber
	}

	if req.Capacity != nil {
		vehicle.Capacity = *req.Capacity
	}

	if req.MaxTask != nil {
		vehicle.MaxTask = req.MaxTask
	}

	if req.DailyWorkTimeStart != nil {
		vehicle.DailyWorkTimeStart = req.DailyWorkTimeStart
	}

	if req.DailyWorkTimeEnd != nil {
		vehicle.DailyWorkTimeEnd = req.DailyWorkTimeEnd
	}

	if req.DailyBreakTimeStart != nil {
		vehicle.DailyBreakTimeStart = req.DailyBreakTimeStart
	}

	if req.DailyBreakTimeEnd != nil {
		vehicle.DailyBreakTimeEnd = req.DailyBreakTimeEnd
	}

	if req.StartLat != nil {
		vehicle.StartLat = req.StartLat
	}

	if req.StartLon != nil {
		vehicle.StartLon = req.StartLon
	}

	if req.EndLat != nil {
		vehicle.EndLat = req.EndLat
	}

	if req.EndLon != nil {
		vehicle.EndLon = req.EndLon
	}

	if err := validators.ValidateVehicle(vehicle); err != nil {
		return nil, err
	}

	if err := s.repo.Update(vehicle); err != nil {
		return nil, err
	}

	if req.TagSkills != nil {

		if err := s.repo.DeleteSkills(
			vehicle.ID,
		); err != nil {
			return nil, err
		}

		var skills []models.VehicleTagSkill

		for _, tag := range *req.TagSkills {

			var tagSkillID uint

			if tag.ID != nil {

				tagSkillID = *tag.ID

			} else {

				existing, err :=
					s.repo.FindTagSkillByName(
						tag.Name,
					)

				if err == nil {

					tagSkillID = existing.ID

				} else {

					newTag := models.TagSkill{
						Name:  tag.Name,
						Color: tag.Color,
					}

					if err := s.repo.CreateTagSkill(
						&newTag,
					); err != nil {
						return nil, err
					}

					tagSkillID = newTag.ID
				}
			}

			skills = append(
				skills,
				models.VehicleTagSkill{
					VehicleID:  vehicle.ID,
					TagSkillID: tagSkillID,
				},
			)
		}

		if len(skills) > 0 {

			if err := s.repo.CreateSkills(
				skills,
			); err != nil {
				return nil, err
			}
		}
	}

	updatedVehicle, err := s.repo.FindByIDWithSkills(
		vehicle.ID,
	)

	if err != nil {
		return nil, err
	}

	responseData, err :=
		s.mapVehicleResponse(
			updatedVehicle,
		)

	if err != nil {
		return nil, err
	}

	return &responseData, nil
}

func (s *VehicleService) Delete(
	req dto.DeleteVehicle,
) error {

	for _, id := range req.ID {

		vehicle, err := s.repo.FindByIDAndPlan(
			id,
			req.PlanID,
		)

		if err != nil {
			return err
		}

		if err := s.repo.DeleteSkills(
			vehicle.ID,
		); err != nil {
			return err
		}

		if err := s.repo.Delete(
			vehicle.ID,
		); err != nil {
			return err
		}
	}

	return nil
}

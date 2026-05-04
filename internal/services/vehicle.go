package services

import (
	"ROP_Backend/internal/dto"
	"ROP_Backend/internal/models"
	"ROP_Backend/internal/repository"
	"ROP_Backend/internal/validators"

	"gorm.io/gorm"
)

type VehicleService struct {
	repo *repository.VehicleRepository
}

func NewVehicleService(db *gorm.DB) *VehicleService {
	repo := repository.NewVehicleRepository(db)
	return &VehicleService{repo: repo}
}

func (s *VehicleService) GroupCreate(userID uint, req dto.GroupCreateVehicle) error {
	var vehicles []models.Vehicle

	for _, v := range req.Vehicles {
		if err := validators.ValidateVehicle(v); err != nil {
			return err
		}

		vehicle := models.Vehicle{
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
		vehicles = append(vehicles, vehicle)
	}
	return s.repo.Create(vehicles)
}

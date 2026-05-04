package repository

import (
	"ROP_Backend/internal/models"

	"gorm.io/gorm"
)

type VehicleRepository struct {
	db *gorm.DB
}

func NewVehicleRepository(db *gorm.DB) *VehicleRepository {
	return &VehicleRepository{db: db}
}

func (r *VehicleRepository) Create(vehicles []models.Vehicle) error {
	return r.db.Create(&vehicles).Error
}

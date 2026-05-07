package validators

import (
	dto "ROP_Backend/internal/dto/request"
	"errors"
)

func ValidateVehicleSkill(v dto.VehicleSkill) error {

	if v.VehicleID == 0 {
		return errors.New("vehicle_id required")
	}

	return nil
}

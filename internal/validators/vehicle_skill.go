package validators

import (
	"ROP_Backend/internal/dto"
	"errors"
)

func ValidateVehicleSkill(v dto.VehicleSkill) error {

	if v.VehicleID == 0 {
		return errors.New("vehicle_id required")
	}

	return nil
}

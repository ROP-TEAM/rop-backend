package validators

import (
	"ROP_Backend/internal/models"
	"errors"
)

func ValidateOrder(o *models.Order) error {

	if o.Name == "" {
		return errors.New("name required")
	}

	if o.Capacity <= 0 {
		return errors.New("capacity must be > 0")
	}

	if o.Type != 0 && o.Type != 1 {
		return errors.New("type must be 0 or 1")
	}

	if o.Priority < 0 || o.Priority > 3 {
		return errors.New("priority must be 0-3")
	}

	if o.TimeWindowStart != nil && o.TimeWindowEnd != nil {
		if *o.TimeWindowStart >= *o.TimeWindowEnd {
			return errors.New("invalid time window")
		}
	}

	if o.DesLatitude == nil || o.DesLongitude == nil {
		return errors.New("destination required")
	}

	if *o.DesLatitude < -90 || *o.DesLatitude > 90 {
		return errors.New("invalid latitude")
	}

	if *o.DesLongitude < -180 || *o.DesLongitude > 180 {
		return errors.New("invalid longitude")
	}
	return nil
}

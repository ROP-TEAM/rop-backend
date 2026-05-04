package validators

import (
	"ROP_Backend/internal/dto"
	"errors"
)

func ValidateVehicle(v dto.CreateVehicle) error {

	if v.PlateNumber == "" {
		return errors.New("plate number required")
	}

	if v.Capacity <= 0 {
		return errors.New("capacity must be > 0")
	}

	if v.MaxTask != nil && *v.MaxTask <= 0 {
		return errors.New("max_task must be > 0")
	}

	if v.DailyWorkTimeStart != nil && v.DailyWorkTimeEnd != nil {

		if *v.DailyWorkTimeStart >= *v.DailyWorkTimeEnd {
			return errors.New("invalid work time (start >= end)")
		}

		if *v.DailyWorkTimeStart < 0 {
			return errors.New("work time out of range")
		}
	}

	if v.DailyBreakTimeStart != nil && v.DailyBreakTimeEnd != nil {

		if *v.DailyBreakTimeStart >= *v.DailyBreakTimeEnd {
			return errors.New("invalid break time")
		}

		if v.DailyWorkTimeStart != nil && v.DailyWorkTimeEnd != nil {
			if *v.DailyBreakTimeStart < *v.DailyWorkTimeStart ||
				*v.DailyBreakTimeEnd > *v.DailyWorkTimeEnd {
				return errors.New("break must be within work time")
			}
		}

		if (*v.DailyBreakTimeEnd - *v.DailyBreakTimeStart) > 180 {
			return errors.New("break too long")
		}
	}

	if (v.StartLat != nil && v.StartLon == nil) ||
		(v.StartLat == nil && v.StartLon != nil) {
		return errors.New("start_lat and start_lon must both be provided")
	}

	if (v.EndLat != nil && v.EndLon == nil) ||
		(v.EndLat == nil && v.EndLon != nil) {
		return errors.New("end_lat and end_lon must both be provided")
	}

	if v.StartLat == nil && v.EndLat != nil {
		return errors.New("cannot have end location without start location")
	}

	return nil
}

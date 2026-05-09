package validators

import (
	"ROP_Backend/internal/models"
	"errors"
)

func ValidateVehicleCreate(v *models.Vehicle) error {

	if v.PlateNumber == "" {
		return errors.New("plate number required")
	}

	if v.Capacity <= 0 {
		return errors.New("capacity must be > 0")
	}

	if v.MaxTask != nil && *v.MaxTask < 0 {
		return errors.New("max_task must be >= 0 (use 0 for unlimited)")
	}

	if v.DailyWorkTimeStart == nil || v.DailyWorkTimeEnd == nil {
		return errors.New("work start/end time required")
	}

	if v.DailyBreakTimeStart == nil || v.DailyBreakTimeEnd == nil {
		return errors.New("break start/end time required")
	}

	if *v.DailyWorkTimeStart >= *v.DailyWorkTimeEnd {
		return errors.New("invalid work time (start >= end)")
	}

	if *v.DailyWorkTimeStart < 0 ||
		*v.DailyWorkTimeEnd < 0 {
		return errors.New("work time out of range")
	}

	if *v.DailyBreakTimeStart >= *v.DailyBreakTimeEnd {
		return errors.New("invalid break time")
	}

	if *v.DailyBreakTimeStart < *v.DailyWorkTimeStart ||
		*v.DailyBreakTimeEnd > *v.DailyWorkTimeEnd {
		return errors.New("break must be within work time")
	}

	if (*v.DailyBreakTimeEnd - *v.DailyBreakTimeStart) > 180 {
		return errors.New("break too long")
	}

	if (v.StartLat != nil && v.StartLon == nil) ||
		(v.StartLat == nil && v.StartLon != nil) {
		return errors.New("start_lat and start_lon must both be provided")
	}

	if (v.EndLat != nil && v.EndLon == nil) ||
		(v.EndLat == nil && v.EndLon != nil) {
		return errors.New("end_lat and end_lon must both be provided")
	}

	if v.StartLat != nil && (*v.StartLat < -90 || *v.StartLat > 90) {
		return errors.New("start_lat must be between -90 and 90")
	}

	if v.StartLon != nil && (*v.StartLon < -180 || *v.StartLon > 180) {
		return errors.New("start_lon must be between -180 and 180")
	}

	if v.EndLat != nil && (*v.EndLat < -90 || *v.EndLat > 90) {
		return errors.New("end_lat must be between -90 and 90")
	}

	if v.EndLon != nil && (*v.EndLon < -180 || *v.EndLon > 180) {
		return errors.New("end_lon must be between -180 and 180")
	}

	if v.StartLat == nil && v.EndLat != nil {
		return errors.New("cannot have end location without start location")
	}

	return nil
}

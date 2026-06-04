package validators

import (
	"errors"
	"time"
)

const (
	PlanDateLayout = "2006-01-02 15:04:05.999999-07"
)

func ValidatePlanDate(dateStr string) (time.Time, error) {
	t, err := time.Parse(PlanDateLayout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("plan_date must be in format: \"2026-05-05 15:44:09.523069+00\"")
	}
	return t, nil
}

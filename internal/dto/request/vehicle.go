package dto

type CreateVehicle struct {
	PlateNumber string  `json:"plate_number"`
	Model       string  `json:"model"`
	Name        string  `json:"name"`
	Capacity    float64 `json:"capacity"`
	MaxTask     *int    `json:"max_task"`

	DailyWorkTimeStart  *int `json:"daily_work_time_start"`
	DailyWorkTimeEnd    *int `json:"daily_work_time_end"`
	DailyBreakTimeStart *int `json:"daily_break_time_start"`
	DailyBreakTimeEnd   *int `json:"daily_break_time_end"`

	StartLat *float64 `json:"start_latitude"`
	StartLon *float64 `json:"start_longitude"`
	EndLat   *float64 `json:"end_latitude"`
	EndLon   *float64 `json:"end_longitude"`
}

type GroupCreateVehicle struct {
	PlanID   string          `json:"plan_id"`
	Vehicles []CreateVehicle `json:"vehicles"`
}

package dto

type OptimizeRequest struct {
	DepotLat float64 `json:"depot_lat" validate:"required"`
	DepotLon float64 `json:"depot_lon" validate:"required"`

	Vehicles []OptimizeVehicle `json:"vehicles" validate:"required,min=1"`
	Orders   []OptimizeOrder   `json:"orders" validate:"required,min=1"`

	EnableALNS      bool   `json:"enable_alns"`
	EnableMultiTrip bool   `json:"enable_multi_trip"`
	ReloadMin       int    `json:"reload_min"`    // minutes; 0 = server default (30)
	TimeLimitMS     int    `json:"time_limit_ms"` // 0 = server default (5000)
	Seed            uint32 `json:"seed"`          // 0 = non-deterministic
}

type OptimizeVehicle struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	PlateNumber string `json:"plate_number"`

	Capacity int `json:"capacity"`
	MaxTask  int `json:"max_task"`

	DailyWorkTimeStart int `json:"daily_work_time_start"`
	DailyWorkTimeEnd   int `json:"daily_work_time_end"`

	DailyBreakTimeStart int `json:"daily_break_time_start"`
	DailyBreakTimeEnd   int `json:"daily_break_time_end"`

	StartLatitude  float64 `json:"start_latitude"`
	StartLongitude float64 `json:"start_longitude"`

	EndLatitude  float64 `json:"end_latitude"`
	EndLongitude float64 `json:"end_longitude"`

	Skills []CreateTagSkill `json:"skills"`
}

type OptimizeOrder struct {
	Name string `json:"name"`

	Type int `json:"type"`

	Capacity int `json:"capacity"`

	DesLatitude  float64 `json:"des_latitude"`
	DesLongitude float64 `json:"des_longitude"`

	TimeWindowStart int `json:"time_window_start"`
	TimeWindowEnd   int `json:"time_window_end"`

	ServiceTime int `json:"service_time"`

	Priority int `json:"priority"`

	Skills []CreateTagSkill `json:"skills"`
}

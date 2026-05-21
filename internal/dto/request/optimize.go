package dto

type OptimizeRequest struct {
	DepotLat float64 `json:"depotLat" validate:"required"`
	DepotLon float64 `json:"depotLon" validate:"required"`

	Vehicles []OptimizeVehicle `json:"vehicles" validate:"required,min=1"`
	Orders   []OptimizeOrder   `json:"orders" validate:"required,min=1"`

	EnableALNS      bool   `json:"enableAlns"`
	EnableMultiTrip bool   `json:"enableMultiTrip"`
	ReloadMin       int    `json:"reloadMin"`   // minutes; 0 = server default (30)
	TimeLimitMS     int    `json:"timeLimitMS"` // 0 = server default (5000)
	Seed            uint32 `json:"seed"`        // 0 = non-deterministic
}

type OptimizeVehicle struct {
	Name        string `json:"name"`
	Model       string `json:"model"`
	PlateNumber string `json:"plateNumber"`

	Capacity int `json:"capacity"`
	MaxTask  int `json:"maxTask"`

	DailyWorkTimeStart int `json:"workTimeStart"`
	DailyWorkTimeEnd   int `json:"workTimeEnd"`

	DailyBreakTimeStart int `json:"breakTimeStart"`
	DailyBreakTimeEnd   int `json:"breakTimeEnd"`

	StartLatitude  float64 `json:"startLatitude"`
	StartLongitude float64 `json:"startLongitude"`

	EndLatitude  float64 `json:"endLatitude"`
	EndLongitude float64 `json:"endLongitude"`

	Skills []CreateTagSkill `json:"skills"`
}

type OptimizeOrder struct {
	Name string `json:"name"`

	Type int `json:"type"`

	Capacity int `json:"capacity"`

	DesLatitude  float64 `json:"desLatitude"`
	DesLongitude float64 `json:"desLongitude"`

	TimeWindowStart int `json:"timeWindowStart"`
	TimeWindowEnd   int `json:"timeWindowEnd"`

	ServiceTime int `json:"serviceTime"`

	Priority int `json:"priority"`

	Skill string `json:"skill"`
}

package dto

type OptimizeRequest struct {
	DepotLat float64 `json:"depotLat" validate:"gte=-90,lte=90"`
	DepotLon float64 `json:"depotLon" validate:"gte=-180,lte=180"`

	Vehicles []OptimizeVehicle `json:"vehicles" validate:"required,min=1,dive"`
	Orders   []OptimizeOrder   `json:"orders" validate:"required,min=1,dive"`

	EnableALNS      bool   `json:"enableAlns"`
	EnableMultiTrip bool   `json:"enableMultiTrip"`
	ReloadMin       int    `json:"reloadMin"`   // minutes; 0 = server default (30)
	TimeLimitMS     int    `json:"timeLimitMS"` // 0 = server default (5000)
	Seed            uint32 `json:"seed"`        // 0 = non-deterministic
}

type OptimizeVehicle struct {
	Name        string `json:"name" validate:"required"`
	Model       string `json:"model"`
	PlateNumber string `json:"plateNumber"`

	Capacity int `json:"capacity" validate:"gt=0"`
	MaxTask  int `json:"maxTask" validate:"gte=0"`

	DailyWorkTimeStart int `json:"workTimeStart" validate:"gte=0,lte=1439"`
	DailyWorkTimeEnd   int `json:"workTimeEnd" validate:"gte=0,lte=1439"`

	DailyBreakTimeStart int `json:"breakTimeStart" validate:"gte=0,lte=1439"`
	DailyBreakTimeEnd   int `json:"breakTimeEnd" validate:"gte=0,lte=1439"`

	StartLatitude  float64 `json:"startLatitude" validate:"gte=-90,lte=90"`
	StartLongitude float64 `json:"startLongitude" validate:"gte=-180,lte=180"`

	EndLatitude  float64 `json:"endLatitude" validate:"gte=-90,lte=90"`
	EndLongitude float64 `json:"endLongitude" validate:"gte=-180,lte=180"`

	Skills []CreateTagSkill `json:"skills"`
}

type OptimizeOrder struct {
	Name string `json:"name" validate:"required"`

	Type int `json:"type" validate:"oneof=0 1"`

	Capacity int `json:"capacity" validate:"gt=0"`

	DesLatitude  float64 `json:"desLatitude" validate:"gte=-90,lte=90"`
	DesLongitude float64 `json:"desLongitude" validate:"gte=-180,lte=180"`

	TimeWindowStart int `json:"timeWindowStart" validate:"gte=0,lte=1439"`
	TimeWindowEnd   int `json:"timeWindowEnd" validate:"gte=0,lte=1439"`

	ServiceTime int `json:"serviceTime" validate:"gte=0"`

	Priority int `json:"priority" validate:"oneof=0 1 2 3"`

	Skill string `json:"skill"`
}

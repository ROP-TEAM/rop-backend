package response

type OptimizeResponse struct {
	Message     string               `json:"message"`
	Routes      []RouteResponse      `json:"routes"`
	Unassigned  []string             `json:"unassigned,omitempty"`
	DropReasons []DropReasonResponse `json:"dropReasons,omitempty"`
	DepotLat    float64              `json:"depotLat"`
	DepotLon    float64              `json:"depotLon"`
}

type DropReasonResponse struct {
	Code   string `json:"code"`
	Detail string `json:"detail"`
}

type VehicleSkill struct {
	Name  string  `json:"name"`
	ID    *int    `json:"id,omitempty"`
	Color *string `json:"color,omitempty"`
}

type RouteResponse struct {
	TotalDistance float64 `json:"totalDistance"`
	TotalDuration float64 `json:"totalDuration"`

	//vehicle identity
	Name           string         `json:"name"`
	Capacity       int            `json:"capacity"`
	WorkTimeStart  int            `json:"workTimeStart"`
	WorkTimeEnd    int            `json:"workTimeEnd"`
	Model          *string        `json:"model,omitempty"`
	PlateNumber    *string        `json:"plateNumber,omitempty"`
	ProfileID      *int           `json:"profile_id,omitempty"`
	BreakTimeStart *int           `json:"breakTimeStart,omitempty"`
	BreakTimeEnd   *int           `json:"breakTimeEnd,omitempty"`
	MaxTask        *int           `json:"maxTask,omitempty"`
	Skills         []VehicleSkill `json:"skills,omitempty"`
	ID             int            `json:"id"`
	Stops          []StopResponse `json:"stops"`
	// TripSizes []int          `json:"tripSizes,omitempty"`
}

type StopResponse struct {
	OrderName string `json:"orderName"`

	ArrivalMin int `json:"arrivalMin"`
	// DepartMin  int `json:"departMin"`

	DistanceFromPrevious float64 `json:"distanceFromPrevious"` // meters
	DurationFromPrevious int     `json:"durationFromPrevious"` // minutes

	//order identity
	Capacity        int     `json:"capacity"`
	TimeWindowStart int     `json:"timeWindowStart"`
	TimeWindowEnd   int     `json:"timeWindowEnd"`
	DesLatitude     float64 `json:"desLatitude"`
	DesLongitude    float64 `json:"desLongitude"`
	ServiceTime     int     `json:"serviceTime"`
	Type            int     `json:"type"`
	Priority        int     `json:"priority"`
	Note            *string `json:"note,omitempty"`
	Skill           *string `json:"skill,omitempty"`
	ID              int     `json:"id"`
}

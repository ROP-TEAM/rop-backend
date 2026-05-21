package response

type OptimizeResponse struct {
	Status      string               `json:"status"`
	Routes      []RouteResponse      `json:"routes"`
	Unassigned  []string             `json:"unassigned,omitempty"`
	DropReasons []DropReasonResponse `json:"dropReasons,omitempty"`
}

type DropReasonResponse struct {
	OrderName string `json:"orderName"`
	Code      string `json:"code"`
	Detail    string `json:"detail"`
}

type RouteResponse struct {
	VehicleName string `json:"vehicleName"`

	TotalDistance float64 `json:"totalDistance"`
	TotalDuration float64 `json:"totalDuration"`

	Stops     []StopResponse `json:"stops"`
	TripSizes []int          `json:"tripSizes,omitempty"`
}

type StopResponse struct {
	OrderName string `json:"orderName"`

	ArrivalMin int `json:"arrivalMin"`
	DepartMin  int `json:"departMin"`
}

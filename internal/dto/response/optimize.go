package response

type OptimizeResponse struct {
	Status      string               `json:"status"`
	Routes      []RouteResponse      `json:"routes"`
	Unassigned  []string             `json:"unassigned,omitempty"`
	DropReasons []DropReasonResponse `json:"drop_reasons,omitempty"`
}

type DropReasonResponse struct {
	OrderName string `json:"order_name"`
	Code      string `json:"code"`
	Detail    string `json:"detail"`
}

type RouteResponse struct {
	VehicleName string `json:"vehicle_name"`

	TotalDistance float64 `json:"total_distance"`
	TotalDuration float64 `json:"total_duration"`

	Stops []StopResponse `json:"stops"`
}

type StopResponse struct {
	OrderName string `json:"order_name"`

	ArrivalMin int `json:"arrival_min"`
	DepartMin  int `json:"depart_min"`
}

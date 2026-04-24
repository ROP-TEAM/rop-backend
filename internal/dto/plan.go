package dto

type PlanRequest struct {
	Vehicles []VehiclesRequest `json:"vehicle"`
	Orders   []OrderRequest    `json:"orders"`
}

type VehiclesRequest struct {
	ID       int      `json:"id"`
	Model    string   `json:"model"`
	Display  string   `json:"display"`
	Capacity int      `json:"capacity"`
	Skills   []string `json:"skills"`
	MaxTask  int      `json:"maxTask"`

	WorkTime  TimeRange `json:"workTime"`
	BreakTime TimeRange `json:"breakTime"`

	StartLocation Location `json:"startLocation"`
	EndLocation   Location `json:"endLocation"`
}

type OrderRequest struct {
	ID       int      `json:"id"`
	Name     string   `json:"name"`
	Capacity int      `json:"capacity"`
	Skills   []string `json:"skills"`

	Location Location `json:"location"`

	TimeWindow     TimeRange `json:"timeWindow"`
	DeliveryWindow TimeRange `json:"deliveryWindow"`

	Type     string `json:"type"`
	Priority string `json:"priority"`

	ServiceTime       int `json:"serviceTime"`
	AssignToVehicleID int `json:"assignToVehicleID"`
}

type TimeRange struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

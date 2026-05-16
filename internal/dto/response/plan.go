package response

import "time"

type StopResponse struct {
	ID             uint          `json:"id"`
	SequenceNumber int           `json:"sequence_number"`
	OrderID        uint          `json:"order_id"`
	ArrivalMin     int           `json:"arrival_min"`
	DepartMin      int           `json:"depart_min"`
}

type RouteResponse struct {
	ID            uint           `json:"id"`
	VehicleID     uint           `json:"vehicle_id"`
	TotalDistance *float64       `json:"total_distance"`
	TotalTime     *int           `json:"total_time"`
	Stops         []StopResponse `json:"stops"`
}

type PlanResponse struct {
	ID        string          `json:"id"`
	Name      string          `json:"name"`
	PlanDate  time.Time       `json:"plan_date"`
	Status    string          `json:"status"`
	DepotLat  float64         `json:"depot_lat"`
	DepotLon  float64         `json:"depot_lon"`
	Routes    []RouteResponse `json:"routes"`
	CreatedAt time.Time       `json:"created_at"`
}

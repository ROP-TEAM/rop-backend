package dto

type CreatePlan struct {
	Name     string  `json:"name"`
	PlanDate string  `json:"plan_date"` // "2006-01-02"
	DepotLat float64 `json:"depot_lat"`
	DepotLon float64 `json:"depot_lon"`
}

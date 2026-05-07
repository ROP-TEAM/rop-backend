package response

type TagSkillResponse struct {
	ID           uint   `json:"id"`
	Name         string `json:"name"`
	Color        string `json:"color"`
	OrderCount   int64  `json:"order_count"`
	VehicleCount int64  `json:"vehicle_count"`
}

package response

type VehicleResponse struct {
	VehicleID   uint    `json:"vehicle_id"`
	ProfileID   *int    `json:"profile_id"`
	PlateNumber string  `json:"plate_number"`
	Model       string  `json:"model"`
	Name        string  `json:"name"`
	Capacity    float64 `json:"capacity"`
	MaxTask     *int    `json:"max_task"`

	TagSkillID []uint `json:"tag_skill_id"`
}

type VehicleGroupResponse struct {
	Message string            `json:"message"`
	Data    []VehicleResponse `json:"data"`
}

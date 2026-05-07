package dto

type VehicleSkill struct {
	VehicleID  uint   `json:"vehicle_id"`
	TagSkillID []uint `json:"tag_skill_id"`
}

type GroupVehicleSkill struct {
	PlanID   string         `json:"plan_id"`
	Vehicles []VehicleSkill `json:"vehicles"`
}

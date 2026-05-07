package dto

type CreateOrder struct {
	Name string `json:"name"`
	Note string `json:"note"`

	Type int `json:"type"` // 0=delivery,1=pickup

	Capacity float64 `json:"capacity"`

	ServiceTime *int `json:"service_time"`
	Priority    int  `json:"priority"` // 0-3

	TimeWindowStart *int `json:"time_window_start"`
	TimeWindowEnd   *int `json:"time_window_end"`

	DesLatitude  float64 `json:"des_latitude"`
	DesLongitude float64 `json:"des_longitude"`

	TagSkillID []uint `json:"tag_skill_id"`
}

type GroupCreateOrder struct {
	PlanID string        `json:"plan_id"`
	Orders []CreateOrder `json:"orders"`
}

package response

type OrderResponse struct {
	Name string `json:"name"`
	Note string `json:"note"`

	Type int `json:"type"` // 0=delivery,1=pickup

	Capacity float64 `json:"capacity"`

	ServiceTime *int `json:"service_time"`
	Priority    int  `json:"priority"` // 0-3

	TagSkillID []uint `json:"tag_skill_id"`
}

type OrderGroupResponse struct {
	Message string          `json:"message"`
	Data    []OrderResponse `json:"data"`
}

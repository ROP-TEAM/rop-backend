package dto

type CreateTagSkill struct {
	ID    *uint  `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type GroupCreateTagSkill struct {
	PlanID string           `json:"plan_id"`
	Skills []CreateTagSkill `json:"skills"`
}

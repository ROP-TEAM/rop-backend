package dto

type CreatePlanRequest struct {
	Name string `gorm:"type:varchar(100) default:untitle" json:"name" example:"untitle"`
}

type CreatePlanResponse struct {
	Name string `gorm:"type:varchar(100) default:untitle" json:"name" example:"what"`
}

package dto

import (
	"ROP_Backend/internal/models"
	"time"
)

type CreatePlanRequest struct {
	Name     *string   `json:"name" example:"untangle"`
	PlanDate time.Time `json:"plan_date" example:"12-12-12"`
}

func (r *CreatePlanRequest) SetDefaults() {
	if *r.Name == "" {
		*r.Name = "untitle"
	}
	if r.PlanDate.IsZero() {
		r.PlanDate = time.Now()
	}
}

type CreatePlanResponse struct {
	ID        string    `json:"id" example:"whatthehell"`
	Name      string    `json:"name" example:"what"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdatePlanNameByIDRequest struct {
	Name string `json:"name" example:"united in grief"`
}

func (r *UpdatePlanNameByIDRequest) SetDefaults() {
	if r.Name == "" {
		r.Name = "untitle"
	}
}

type UpdatePlanNameByIDResponse struct {
	Name      string    `json:"name" example:"united in grief"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeletePlanByIDRequest struct {
	ID string `json:"id" example:"whatthehelpyouaskfor"`
}

type DeletePlanByIDResponse struct {
}

type GetPlanByIDRequest struct {
	ID string `json:"id" example:"whatthehelpyouaskfor"`
}

type GetPlanByIDResponse struct {
	ID        string    `json:"id" example:"whatthehell"`
	Name      string    `json:"name" example:"what"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Vehicles  []models.Vehicle
}

type GetPlansRequest struct {
}

type GetPlansResponse struct {
	plans []GetPlanByIDResponse
}

// zignimaaa เพราว่าต้องไปนั่งduplicate data แต่ละอันทีละtable
// ถามว่าเอามาทั้งorder vehicle
type DuplicatePlanByIDRequest struct {
}

type DuplicatePlanByIDResponse struct {
}

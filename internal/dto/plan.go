package dto

import (
	"ROP_Backend/internal/models"
	"strings"
	"time"
)

const (
	DefautlPlanName = "untitle"
)

func isSpaceName(a *string) bool {
	return a == nil || strings.TrimSpace(*a) == ""
}

type CreatePlanRequest struct {
	Name     *string   `json:"name" example:"untangle"`
	PlanDate time.Time `json:"plan_date" example:"2026-05-05 15:44:09.523069+00"`
}

func (r *CreatePlanRequest) SetDefaults() {
	if isSpaceName(r.Name) {
		*r.Name = DefautlPlanName
	}
	if r.PlanDate.IsZero() {
		r.PlanDate = time.Now()
	}
}

type CreatePlanResponse struct {
	ID        string    `json:"id" example:"whatthehell"`
	Name      string    `json:"name" example:"what"`
	CreatedAt time.Time `json:"created_at" example:"2026-05-05 15:44:09.523069+00"`
}

type UpdatePlanNameByIDRequest struct {
	Name *string `json:"name" example:"united in grief"`
}

func (r *UpdatePlanNameByIDRequest) SetDefaults() {
	if isSpaceName(r.Name) {
		*r.Name = DefautlPlanName
	}
}

type UpdatePlanNameByIDResponse struct {
	Name      string    `json:"name" example:"united in grief"`
	UpdatedAt time.Time `json:"updated_at"`
}

type DeletePlanByIDRequest struct {
}

type DeletePlanByIDResponse struct {
	Message string `example:"plan deleted"`
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

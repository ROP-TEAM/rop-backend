package dto

import (
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

// type GetPlanByIDResponse struct {
// 	ID        string    `json:"id" example:"whatthehell"`
// 	Name      string    `json:"name" example:"what"`
// 	CreatedAt time.Time `json:"created_at"`
// 	UpdatedAt time.Time `json:"updated_at"`
// 	Vehicles  []models.Vehicle
// }

type GetPlanMetaDataRequest struct {
}

type GetPlanMetaDataResponse struct {
	ID        string    `json:"id" example:"whatthehell"`
	Name      string    `json:"name" example:"what"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetPlansRequest struct {
}

type TagSkillDetails struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Color string `json:"color"`
}

type TagSkillID struct {
	ID uint `json:"id"`
}

type VehicleDetails struct {
	ID                  uint         `json:"id"`
	ProfileID           *int         `json:"profile_id"`
	PlateNumber         string       `json:"plate_number"`
	Model               string       `json:"model"`
	Name                string       `json:"name"`
	Capacity            float64      `json:"capacity"`
	MaxTask             *int         `json:"max_task"`
	DailyBreakTimeStart *int         `json:"daily_break_time_start"`
	DailyBreakTimeEnd   *int         `json:"daily_break_time_end"`
	DailyWorkTimeStart  *int         `json:"daily_work_time_start"`
	DailyWorkTimeEnd    *int         `json:"daily_work_time_end"`
	StartLat            *float64     `json:"start_latitude"`
	StartLon            *float64     `json:"start_longitude"`
	EndLat              *float64     `json:"end_latitude"`
	EndLon              *float64     `json:"end_longitude"`
	Skills              []TagSkillID `json:"skills"`
}

type OrderDetails struct {
	ID              uint         `json:"id"`
	Name            string       `json:"name"`
	Note            string       `json:"note"`
	Type            int          `json:"type"`
	Capacity        float64      `json:"capacity"`
	ServiceTime     *int         `json:"service_time"`
	Priority        int          `json:"priority"`
	TimeWindowStart *int         `json:"time_window_start"`
	TimeWindowEnd   *int         `json:"time_window_end"`
	DesLatitude     *float64     `json:"des_latitude"`
	DesLongitude    *float64     `json:"des_longitude"`
	Skills          []TagSkillID `json:"skills"`
}

type StopDetails struct {
	ID             uint `json:"id"`
	SequenceNumber int  `json:"sequence_number"`
	RouteID        uint `json:"route_id"`
	OrderID        uint `json:"order_id"`
}

type RouteDetails struct {
	ID            uint          `json:"id"`
	VehicleID     uint          `json:"vehicle_id"`
	TotalDistance *float64      `json:"total_distance"`
	TotalTime     *int          `json:"total_time"`
	Stops         []StopDetails `json:"stops"`
}

// type TagSkillDetails struct {
// 	ID    uint   `json:"id"`
// 	Name  string `json:"name"`
// 	Color string `json:"color"`
// }

type GetPlansDetailsResponse struct {
	ID        string            `json:"id" example:"whatthehell"`
	Name      string            `json:"name" example:"what"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
	Vehicles  []VehicleDetails  `json:"vehicles"`
	Orders    []OrderDetails    `json:"orders"`
	Routes    []RouteDetails    `json:"routes"`
	TagSkills []TagSkillDetails `json:"tag_skills"`
}

// zignimaaa เพราว่าต้องไปนั่งduplicate data แต่ละอันทีละtable
// ถามว่าเอามาทั้งorder vehicle
type DuplicatePlanByIDRequest struct {
}

type DuplicatePlanByIDResponse struct {
	ID   string `json:"id" example:"whatthehell"`
	Name string `json:"" example:"vhkp9^,lvo0y[0bh'0dsojvp]"`
}

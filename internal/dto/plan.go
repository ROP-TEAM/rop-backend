package dto

import (
	"time"
)

type CreatePlanRequest struct {
	Name     *string `json:"name" example:"untangle"`
	PlanDate *string `json:"plan_date" example:"2026-05-05 15:44:09.523069+00"` // CHANGED: time.Time -> *string
}

type CreatePlanResponse struct {
	ID        string    `json:"id" example:"whatthehell"`
	Name      string    `json:"name" example:"what"`
	CreatedAt time.Time `json:"created_at" example:"2026-05-05 15:44:09.523069+00"`
}

type UpdatePlanNameByIDRequest struct {
	Name *string `json:"name" example:"united in grief"`
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

type GetPlanMetaDataResponse struct {
	ID        string    `json:"id" example:"whatthehell"`
	Name      string    `json:"name" example:"what"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type GetPlansResponse struct {
	Plans []GetPlanMetaDataResponse `json:"plans"`
}

type TagSkillDetails struct {
	ID    uint   `json:"id" example:"1314"`
	Name  string `json:"name" example:"driver must meet buety standard"`
	Color string `json:"color" example:"red"`
}

type TagSkillID struct {
	ID uint `json:"id" example:"1314"`
}

type VehicleDetails struct {
	ID                  uint         `json:"id" example:"1314"`
	ProfileID           *int         `json:"profile_id" example:"1"`
	NumberPlate         string       `json:"number_plate" example:"สกร2019"`
	Model               string       `json:"model" example:"Malah model Y"`
	Name                string       `json:"name" example:"รถคันนี้สีแดง"`
	Capacity            float64      `json:"capacity" example:"1000"`
	MaxTask             *int         `json:"max_task" example:"20"`
	DailyBreakTimeStart *int         `json:"daily_break_time_start" example:"720"`
	DailyBreakTimeEnd   *int         `json:"daily_break_time_end" example:"780"`
	DailyWorkTimeStart  *int         `json:"daily_work_time_start" example:"480"`
	DailyWorkTimeEnd    *int         `json:"daily_work_time_end" example:"960"`
	StartLat            *float64     `json:"start_latitude" example:"1.3288"`
	StartLon            *float64     `json:"start_longitude" example:"8.8231"`
	EndLat              *float64     `json:"end_latitude" example:"1.3288"`
	EndLon              *float64     `json:"end_longitude" example:"8.8231"`
	Skills              []TagSkillID `json:"skills"`
}

type OrderDetails struct {
	ID              uint         `json:"id" example:"720"`
	Name            string       `json:"name" example:"red bicycle"`
	Note            string       `json:"note" example:"owned by P sek"`
	Type            int          `json:"type" example:"1"`
	Capacity        float64      `json:"capacity" example:"77"`
	ServiceTime     *int         `json:"service_time" example:"7"`
	Priority        int          `json:"priority" example:"1"`
	TimeWindowStart *int         `json:"time_window_start" example:"720"`
	TimeWindowEnd   *int         `json:"time_window_end" example:"1000"`
	DesLatitude     *float64     `json:"des_latitude" example:"1.992"`
	DesLongitude    *float64     `json:"des_longitude" example:"2.991"`
	Skills          []TagSkillID `json:"skills"`
}

type StopDetails struct {
	ID             uint `json:"id" example:"7"`
	SequenceNumber int  `json:"sequence_number" example:"77"`
	RouteID        uint `json:"route_id" example:"72"`
	OrderID        uint `json:"order_id" example:"720"`
}

type RouteDetails struct {
	ID            uint          `json:"id" example:"720"`
	VehicleID     uint          `json:"vehicle_id" example:"1314"`
	TotalDistance *float64      `json:"total_distance" example:"720"`
	TotalTime     *int          `json:"total_time" example:"720"`
	Stops         []StopDetails `json:"stops"`
}

// type TagSkillDetails struct {
// 	ID    uint   `json:"id"`
// 	Name  string `json:"name"`
// 	Color string `json:"color"`
// }

type GetPlanDetailsResponse struct {
	ID        string            `json:"id" example:"whatthehell"`
	Name      string            `json:"name" example:"what"`
	CreatedAt time.Time         `json:"created_at" example:"2026-05-05 15:44:09.523069+00"`
	UpdatedAt time.Time         `json:"updated_at" example:"2026-05-05 15:44:09.523069+00"`
	Vehicles  []VehicleDetails  `json:"vehicles"`
	Orders    []OrderDetails    `json:"orders"`
	Routes    []RouteDetails    `json:"routes"`
	TagSkills []TagSkillDetails `json:"tag_skills"`
}

type DuplicatePlanByIDRequest struct {
}

type DuplicatePlanByIDResponse struct {
	ID   string `json:"id" example:"whatthehell"`
	Name string `json:"" example:"vhkp9^,lvo0y[0bh'0dsojvp]"`
}

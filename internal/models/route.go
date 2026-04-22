package models

import "time"

type Route struct {
	ID     uint   `gorm:"primaryKey" json:"id"`
	CarID  uint   `gorm:"column:route_car_fk" json:"car_id"`
	Car    Car    `gorm:"foreignKey:CarID" json:"car,omitempty"` // GORM Magic: Association
	PlanID string `gorm:"index;type:uuid"`
	Plan   Plan   `gorm:"foreignKey:PlanID"`

	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`

	TotalDistance float64 `json:"total_distance"`
	TotalTime     int     `json:"total_time"`

	Stops []Stop `gorm:"foreignKey:RouteID" json:"stops"`
}

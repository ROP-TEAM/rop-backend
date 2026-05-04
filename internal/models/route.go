package models

type Route struct {
	ID        uint    `gorm:"primaryKey" json:"id"`
	VehicleID uint    `gorm:"column:route_vehicle_fk" json:"vehicle_id"`
	Vehicle   Vehicle `gorm:"foreignKey:VehicleID" json:"vehicle,omitempty"` // GORM Magic: Association
	PlanID    string  `gorm:"index;type:uuid"`
	Plan      Plan    `gorm:"foreignKey:PlanID"`

	// StartTime time.Time `json:"start_time"`
	// EndTime   time.Time `json:"end_time"`

	TotalDistance *float64 `gorm:"type:double precision" json:"total_distance"`
	TotalTime     *int     `json:"total_time"`

	Stops []Stop `gorm:"foreignKey:RouteID" json:"stops"`
}

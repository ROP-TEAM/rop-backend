package models

type Vehicle struct {
	ID          uint   `gorm:"primaryKey"`
	ProfileID   *int   `json:"profile_id"`
	PlateNumber string `gorm:"type:varchar(20);not null;" json:"plate_number"` //เอาuniqueness ออก
	Model       string `gorm:"type:varchar(100)" json:"model"`
	Name        string `gorm:"type:varchar(100)" json:"name"`

	Capacity float64 `gorm:"type:numeric(10,2);not null" json:"capacity"`
	MaxTask  *int    `gorm:"default:20" json:"max_task"`

	DailyBreakTimeStart *int `gorm:"" json:"daily_break_time_start"`
	DailyBreakTimeEnd   *int `gorm:"" json:"daily_break_time_end"`
	DailyWorkTimeStart  *int `gorm:"" json:"daily_work_time_start"`
	DailyWorkTimeEnd    *int `gorm:"" json:"daily_work_time_end"`

	StartLat *float64 `gorm:"type:double precision" json:"start_latitude"`
	StartLon *float64 `gorm:"type:double precision" json:"start_longitude"`
	EndLat   *float64 `gorm:"type:double precision" json:"end_latitude"`
	EndLon   *float64 `gorm:"type:double precision" json:"end_longitude"`

	// many-to-many relatioinship
	Skills []TagSkill `gorm:"many2many:vehicle_tag_skills" json:"skills"`

	PlanID string `gorm:"column:vehicle_plan_fk;type:uuid;index" json:"plan_id"`
	Plan   Plan   `gorm:"foreignKey:PlanID"`
}

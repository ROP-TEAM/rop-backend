package models

import "time"

type Stop struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	SequenceNumber int       `gorm:"not null;uniqueIndex:idx_route_sequence"`
	ArrivalTime    time.Time `json:"arrival_time"`
	RouteID        uint      `gorm:"index;not null" json:"route_id"`
	Route          Route     `gorm:"foreignKey:RouteID" json:"-"`
	OrderID        uint      `gorm:"not null" json:"order_id"`        // หรือใช้ string ถ้า OrderID เป็น UUID
	Order          Order     `gorm:"foreignKey:OrderID" json:"order"` // ไว้ทำ Preload ดึงพิกัด/ชื่อลูกค้า
}

package models

type Stop struct {
	ID             uint  `gorm:"primaryKey" json:"id"`
	SequenceNumber int   `gorm:"uniqueIndex:idx_route_sequence;not null" json:"sequence_number"`
	RouteID        uint  `gorm:"uniqueIndex:idx_route_sequence;not null" json:"route_id"`
	Route          Route `gorm:"foreignKey:RouteID" json:"-"`
	OrderID        uint  `gorm:"not null" json:"order_id"`        // หรือใช้ string ถ้า OrderID เป็น UUID
	Order          Order `gorm:"foreignKey:OrderID" json:"order"` // ไว้ทำ Preload ดึงพิกัด/ชื่อลูกค้า
	ArrivalMin     int   `json:"arrival_min"`
	DepartMin      int   `json:"depart_min"`
}

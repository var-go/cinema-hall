package models

type SeatType string

const (
	SeatTypeStandard   SeatType = "standard"
	SeatTypeVip        SeatType = "vip"
	SeatTypeWheelchair SeatType = "wheelchair"
)

var SeatTypePrices = map[SeatType]int{
	SeatTypeStandard:   300,
	SeatTypeVip:        600,
	SeatTypeWheelchair: 150,
}

type Seat struct {
	Base
	HallID uint     `json:"hall_id" gorm:"not null"`
	Hall   Hall     `json:"-"`
	Number int      `json:"number" gorm:"not null;uniqueIndex:idx_hall_row_number"`
	Row    int      `json:"row" gorm:"not null;uniqueIndex:idx_hall_row_number"`
	Type   SeatType `json:"type" gorm:"default:'standard'"`
}

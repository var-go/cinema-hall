package models

type Hall struct {
	Base
	Number int    `json:"number" gorm:"not null;uniques"`
	Seats  []Seat `json:"seats"`
}

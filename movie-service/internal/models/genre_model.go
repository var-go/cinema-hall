package models

type Genre struct {
	Base
	Name string `json:"name" gorm:"type:varchar(100);not null;uniqueIndex"`
}

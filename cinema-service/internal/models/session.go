package models

import "time"

type SessionStatus string

const (
	SessionStatusScheduled SessionStatus = "scheduled"
	SessionStatusOngoing   SessionStatus = "ongoing"
	SessionStatusFinished  SessionStatus = "finished"
	SessionStatusCancelled SessionStatus = "cancelled"
)

type Session struct {
	Base
	MovieID   uint          `json:"movie_id" gorm:"not null"`
	HallID    uint          `json:"hall_id" gorm:"not null"`
	StartTime time.Time     `json:"start_time" gorm:"not null"`
	EndTime   time.Time     `json:"end_time" gorm:"not null"`
	Status    SessionStatus `json:"status" gorm:"type:varchar(20);default:'scheduled'"`
}

package dto

import "time"

type CreateSessionRequest struct {
	MovieID   uint      `json:"movie_id" binding:"required"`
	HallID    uint      `json:"hall_id" binding:"required"`
	StartTime time.Time `json:"start_time" binding:"required"`
	EndTime   time.Time `json:"end_time" binding:"required,gtfield=StartTime"`
}

type UpdateSessionRequest struct {
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Status    *string    `json:"status,omitempty" binding:"omitempty,oneof=scheduled ongoing finished cancelled"`
}

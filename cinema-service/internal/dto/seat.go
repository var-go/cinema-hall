package dto

import "cinema-service/internal/models"

type CreateSeatRequest struct {
	Row    int             `json:"row" binding:"required,min=1"`
	Number int             `json:"number" binding:"required,min=1"`
	Type   models.SeatType `json:"type" binding:"omitempty,oneof=standard vip wheelchair"`
}

type UpdateSeatRequest struct {
	Row    *int             `json:"row,omitempty" binding:"omitempty,min=1"`
	Number *int             `json:"number,omitempty" binding:"omitempty,min=1"`
	Type   *models.SeatType `json:"type,omitempty" binding:"omitempty,oneof=standard vip wheelchair"`
}

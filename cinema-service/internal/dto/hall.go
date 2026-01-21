package dto

type CreateHallRequest struct {
	Number int `json:"number" binding:"required"`
}

type UpdateHallRequest struct {
	Number *int `json:"number"`
}

package dto

type GenreCreateRequest struct {
	Name string `json:"name" binding:"required"`
}

type GenreUpdateRequest struct {
	Name *string `json:"name"`
}

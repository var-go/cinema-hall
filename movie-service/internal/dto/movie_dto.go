package dto

import "movie-service/internal/constants"

type MovieCreateRequest struct {
	Title       string                `json:"title" binding:"required"`
	Description string                `json:"description" binding:"required"`
	Year        uint                  `json:"year" binding:"required"`
	Duration    uint                  `json:"duration" binding:"required"`
	AgeRating   string                `json:"age_rating" binding:"required"`
	MovieStatus constants.MovieStatus `json:"movie_status" binding:"required"`
}

type MovieUpdateRequest struct {
	Title       *string                `json:"title"`
	Description *string                `json:"description"`
	Year        *uint                  `json:"year"`
	Duration    *uint                  `json:"duration"`
	AgeRating   *string                `json:"age_rating"`
	MovieStatus *constants.MovieStatus `json:"movie_status"`
}

package models

import "movie-service/internal/constants"

type Movie struct {
	Base
	Title       string                `json:"title" gorm:"type:varchar(255);not null"`
	Description string                `json:"description" gorm:"type:text;not null"`
	Year        uint                  `json:"year" gorm:"not null;index"`
	Duration    uint                  `json:"duration" gorm:"not null"`
	AgeRating   string                `json:"age_rating" gorm:"type:varchar(50);not null"`
	MovieStatus constants.MovieStatus `json:"movie_status" gorm:"type:varchar(50);not null"`
	Genres      []Genre               `json:"genres" gorm:"many2many:movie_genres;"`
}

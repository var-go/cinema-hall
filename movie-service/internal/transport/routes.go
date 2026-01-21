package transport

import (
	"log/slog"
	"movie-service/internal/services"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	routes *gin.Engine,
	movieService services.MovieService,
	genreService services.GenreService,
	logger *slog.Logger,
) {
	movieHandler := NewMovieHandler(movieService, logger)
	genreHandler := NewGenreHandler(genreService, logger)

	movieHandler.RegisterRoutes(routes)
	genreHandler.RegisterRoutes(routes)
}

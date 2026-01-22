package repository

import (
	"log/slog"
	"movie-service/internal/constants"
	"movie-service/internal/models"

	"gorm.io/gorm"
)

type MovieRepository interface {
	Create(movie *models.Movie) error

	List() ([]models.Movie, error)

	GetByID(id uint) (*models.Movie, error)

	GetNowShowing() ([]models.Movie, error)

	GetComingSoon() ([]models.Movie, error)

	Update(movie *models.Movie) error

	Delete(id uint) error
}

type gormMovieRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewMovieRepository(db *gorm.DB, logger *slog.Logger) MovieRepository {
	return &gormMovieRepository{
		DB:     db,
		logger: logger,
	}

}

func (r *gormMovieRepository) Create(movie *models.Movie) error {
	if err := r.DB.Create(movie).Error; err != nil {
		r.logger.Error("failed to create movie", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *gormMovieRepository) List() ([]models.Movie, error) {

	var movies []models.Movie

	if err := r.DB.Preload("Genres").Find(&movies).Error; err != nil {
		r.logger.Error("failed to list movies", slog.Any("error", err))
		return nil, err
	}

	return movies, nil

}

func (r *gormMovieRepository) GetByID(id uint) (*models.Movie, error) {

	var movie models.Movie

	if err := r.DB.Preload("Genres").Where("id = ?", id).First(&movie).Error; err != nil {
		r.logger.Error("failed to get movie by id", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}

	return &movie, nil
}

func (r *gormMovieRepository) GetNowShowing() ([]models.Movie, error) {

	var movies []models.Movie

	if err := r.DB.Preload("Genres").Where("movie_status = ?", constants.MovieNowShowing).Find(&movies).Error; err != nil {
		r.logger.Error("failed to get now showing movies", slog.Any("error", err))
		return nil, err
	}

	return movies, nil
}
func (r *gormMovieRepository) GetComingSoon() ([]models.Movie, error) {

	var movies []models.Movie

	if err := r.DB.Preload("Genres").Where("movie_status = ?", constants.MovieComingSoon).Find(&movies).Error; err != nil {
		r.logger.Error("failed to get coming soon movies", slog.Any("error", err))
		return nil, err
	}

	return movies, nil
}

func (r *gormMovieRepository) Update(movie *models.Movie) error {

	if err := r.DB.Model(&models.Movie{}).Where("id = ?", movie.ID).Updates(movie).Error; err != nil {
		r.logger.Error("failed to update movie", slog.Any("id", movie.ID), slog.Any("error", err))
		return err
	}

	if err := r.DB.Model(movie).Association("Genres").Replace(movie.Genres); err != nil {
		r.logger.Error("failed to update movie genres", slog.Any("id", movie.ID), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormMovieRepository) Delete(id uint) error {

	res := r.DB.Delete(&models.Movie{}, id)
	if err := res.Error; err != nil {
		r.logger.Error("failed to delete movie", slog.Any("id", id), slog.Any("error", err))
		return err
	}

	if res.RowsAffected == 0 {
		r.logger.Info("movie not found for delete", slog.Any("id", id))
		return gorm.ErrRecordNotFound
	}

	return nil
}

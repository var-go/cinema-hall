package repository

import (
	"log/slog"
	"movie-service/internal/models"

	"gorm.io/gorm"
)

type GenreRepository interface {
	Create(genre *models.Genre) error

	List() ([]models.Genre, error)

	GetByID(id uint) (*models.Genre, error)

	Update(genre *models.Genre) error

	Delete(id uint) error
}

type gormGenreRepository struct {
	DB     *gorm.DB
	logger *slog.Logger
}

func NewGenreRepository(db *gorm.DB, logger *slog.Logger) GenreRepository {
	return &gormGenreRepository{
		DB:     db,
		logger: logger,
	}

}

func (r *gormGenreRepository) Create(genre *models.Genre) error {
	if err := r.DB.Create(genre).Error; err != nil {
		r.logger.Error("failed to create genre", slog.Any("error", err))
		return err
	}
	return nil
}

func (r *gormGenreRepository) List() ([]models.Genre, error) {

	var genres []models.Genre

	if err := r.DB.Order("id ASC").Find(&genres).Error; err != nil {
		r.logger.Error("failed to list genres", slog.Any("error", err))
		return nil, err
	}

	return genres, nil

}

func (r *gormGenreRepository) GetByID(id uint) (*models.Genre, error) {

	var genre models.Genre

	if err := r.DB.Where("id = ?", id).First(&genre).Error; err != nil {
		r.logger.Error("failed to get genre by id", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}

	return &genre, nil
}

func (r *gormGenreRepository) Update(genre *models.Genre) error {

	if err := r.DB.Model(&models.Genre{}).Where("id = ?", genre.ID).Updates(genre).Error; err != nil {
		r.logger.Error("failed to update genre", slog.Any("id", genre.ID), slog.Any("error", err))
		return err
	}

	return nil
}

func (r *gormGenreRepository) Delete(id uint) error {

	if err := r.DB.Delete(&models.Genre{}, id).Error; err != nil {
		r.logger.Error("failed to delete genre", slog.Any("id", id), slog.Any("error", err))
		return err
	}

	return nil
}

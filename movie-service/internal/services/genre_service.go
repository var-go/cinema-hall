package services

import (
	"log/slog"
	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/internal/repository"
)

type GenreService interface {
	Create(req *dto.GenreCreateRequest) (*models.Genre, error)

	List() ([]models.Genre, error)

	GetByID(id uint) (*models.Genre, error)

	Update(id uint, req *dto.GenreUpdateRequest) (*models.Genre, error)

	Delete(id uint) error
}

type genreService struct {
	repo   repository.GenreRepository
	logger *slog.Logger
}

func NewGenreService(genreRepo repository.GenreRepository, logger *slog.Logger) GenreService {
	return &genreService{
		repo:   genreRepo,
		logger: logger,
	}
}

func (s *genreService) Create(req *dto.GenreCreateRequest) (*models.Genre, error) {

	genre := models.Genre{
		Name: req.Name,
	}

	if err := s.repo.Create(&genre); err != nil {
		s.logger.Error("genre create failed", slog.Any("error", err), slog.String("name", genre.Name))
		return nil, err
	}

	return &genre, nil
}

func (s *genreService) List() ([]models.Genre, error) {

	genres, err := s.repo.List()

	if err != nil {
		s.logger.Error("genre list failed", slog.Any("error", err))
		return nil, err
	}

	return genres, nil
}

func (s *genreService) GetByID(id uint) (*models.Genre, error) {

	genre, err := s.repo.GetByID(id)

	if err != nil {
		s.logger.Error("genre get by id failed", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}

	return genre, nil
}

func (s *genreService) Update(id uint, req *dto.GenreUpdateRequest) (*models.Genre, error) {

	genre, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("genre update failed: get by id", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}

	if req.Name != nil {
		genre.Name = *req.Name
	}

	if err := s.repo.Update(genre); err != nil {
		s.logger.Error("genre update failed: update", slog.Any("id", genre.ID), slog.Any("error", err))
		return nil, err
	}
	return genre, nil
}

func (s *genreService) Delete(id uint) error {

	err := s.repo.Delete(id)

	if err != nil {
		s.logger.Error("genre delete failed", slog.Any("id", id), slog.Any("error", err))
		return err
	}

	return nil
}

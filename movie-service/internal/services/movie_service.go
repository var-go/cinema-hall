package services

import (
	"log/slog"
	"movie-service/internal/dto"
	"movie-service/internal/models"
	"movie-service/internal/repository"
)

type MovieService interface {
	Create(req *dto.MovieCreateRequest) (*models.Movie, error)

	List() ([]models.Movie, error)

	GetByID(id uint) (*models.Movie, error)

	GetNowShowing() ([]models.Movie, error)

	GetComingSoon() ([]models.Movie, error)

	Update(id uint, req *dto.MovieUpdateRequest) (*models.Movie, error)

	Delete(id uint) error
}

type movieService struct {
	repo      repository.MovieRepository
	genreRepo repository.GenreRepository
	logger    *slog.Logger
}

func NewMovieService(movieRepo repository.MovieRepository, genreRepo repository.GenreRepository, logger *slog.Logger) MovieService {
	return &movieService{
		repo:      movieRepo,
		genreRepo: genreRepo,
		logger:    logger,
	}
}

func (s *movieService) Create(req *dto.MovieCreateRequest) (*models.Movie, error) {

	var genres []models.Genre

	for _, genreID := range req.GenresID {
		genre, err := s.genreRepo.GetByID(genreID)
		if err != nil {
			s.logger.Error("movie create failed: get genre by id", slog.Any("error", err))
			return nil, err
		}
		genres = append(genres, *genre)
	}

	movie := models.Movie{
		Title:       req.Title,
		Description: req.Description,
		Year:        req.Year,
		Duration:    req.Duration,
		AgeRating:   req.AgeRating,
		MovieStatus: req.MovieStatus,
		Genres:      genres,
	}

	if err := s.repo.Create(&movie); err != nil {
		s.logger.Error("movie create failed", slog.Any("error", err), slog.String("title", movie.Title))
		return nil, err
	}

	return &movie, nil
}

func (s *movieService) List() ([]models.Movie, error) {

	movies, err := s.repo.List()

	if err != nil {
		s.logger.Error("movie list failed", slog.Any("error", err))
		return nil, err
	}

	return movies, nil
}

func (s *movieService) GetByID(id uint) (*models.Movie, error) {

	movie, err := s.repo.GetByID(id)

	if err != nil {
		s.logger.Error("movie get by id failed", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}

	return movie, nil
}

func (s *movieService) GetNowShowing() ([]models.Movie, error) {

	movies, err := s.repo.GetNowShowing()

	if err != nil {
		s.logger.Error("failed to get now showing movies", slog.Any("error", err))
		return nil, err
	}

	return movies, nil
}

func (s *movieService) GetComingSoon() ([]models.Movie, error) {

	movies, err := s.repo.GetComingSoon()

	if err != nil {
		s.logger.Error("failed to get coming soon movies", slog.Any("error", err))
		return nil, err
	}

	return movies, nil
}

func (s *movieService) Update(id uint, req *dto.MovieUpdateRequest) (*models.Movie, error) {

	movie, err := s.repo.GetByID(id)
	if err != nil {
		s.logger.Error("movie update failed: get by id", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}

	genre, err := s.genreRepo.GetByID(id)
	if err != nil {
		s.logger.Error("movie update failed: get genre by id", slog.Any("id", id), slog.Any("error", err))
		return nil, err
	}
	movie.Genres = []models.Genre{*genre}

	if req.Title != nil {
		movie.Title = *req.Title
	}

	if req.Description != nil {
		movie.Description = *req.Description
	}

	if req.Year != nil {
		movie.Year = *req.Year
	}

	if req.Duration != nil {
		movie.Duration = *req.Duration
	}

	if req.AgeRating != nil {
		movie.AgeRating = *req.AgeRating
	}

	if req.MovieStatus != nil {
		movie.MovieStatus = *req.MovieStatus
	}

	if err := s.repo.Update(movie); err != nil {
		s.logger.Error("movie update failed: update", slog.Any("id", movie.ID), slog.Any("error", err))
		return nil, err
	}
	return movie, nil
}

func (s *movieService) Delete(id uint) error {

	err := s.repo.Delete(id)

	if err != nil {
		s.logger.Error("movie delete failed", slog.Any("id", id), slog.Any("error", err))
		return err
	}

	return nil
}

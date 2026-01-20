package services

import (
	"cinema-service/internal/dto"
	"cinema-service/internal/models"
	"cinema-service/internal/repository"
	"errors"
	"log/slog"
	"time"
)

type SessionService interface {
	Create(req dto.CreateSessionRequest) (*models.Session, error)
	Update(id uint, req dto.UpdateSessionRequest) (*models.Session, error)
	List() ([]models.Session, error)
	GetById(id uint) (*models.Session, error)
	Delete(id uint) error
	ListByMovieID(movieID uint) ([]models.Session, error)
}

type sessionService struct {
	sessionRepo repository.SessionRepository
	hallRepo    repository.HallRepository
	logger      *slog.Logger
}

func NewSessionService(
	sessionRepo repository.SessionRepository,
	hallRepo repository.HallRepository,
	logger *slog.Logger,
) SessionService {
	return &sessionService{
		sessionRepo: sessionRepo,
		hallRepo:    hallRepo,
		logger:      logger,
	}
}

func (s *sessionService) Create(req dto.CreateSessionRequest) (*models.Session, error) {

	if _, err := s.hallRepo.GetById(req.HallID); err != nil {
		s.logger.Warn(
			"hall not found while creating session",
			"hall_id", req.HallID,
			"error", err,
		)
		return nil, err
	}

	if req.StartTime.Before(time.Now()) {
		s.logger.Warn(
			"attempt to create session in the past",
			"hall_id", req.HallID,
			"movie_id", req.MovieID,
			"start_time", req.StartTime,
		)
		return nil, errors.New("start_time must be in the future")
	}

	session := &models.Session{
		MovieID:   req.MovieID,
		HallID:    req.HallID,
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
		Status:    models.SessionStatusScheduled,
	}

	if err := s.sessionRepo.Create(session); err != nil {
		s.logger.Error(
			"failed to create session",
			"hall_id", req.HallID,
			"movie_id", req.MovieID,
			"start_time", req.StartTime,
			"end_time", req.EndTime,
			"err", err,
		)
		return nil, err
	}

	return session, nil
}

func (s *sessionService) Update(id uint, req dto.UpdateSessionRequest) (*models.Session, error) {

	session, err := s.sessionRepo.GetById(id)
	if err != nil {
		s.logger.Warn(
			"session not found",
			"session_id", id,
			"error", err,
		)
		return nil, err
	}

	if req.StartTime != nil {
		session.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		session.EndTime = *req.EndTime
	}

	if session.EndTime.Before(session.StartTime) || session.EndTime.Equal(session.StartTime) {
		s.logger.Warn(
			"invalid session time range",
			"session_id", id,
			"start_time", session.StartTime,
			"end_time", session.EndTime,
		)
		return nil, errors.New("end_time must be after start_time")
	}

	if req.Status != nil {
		session.Status = models.SessionStatus(*req.Status)
	}

	if err := s.sessionRepo.Update(id, session); err != nil {
		s.logger.Error(
			"failed to update session",
			"session_id", id,
			"hall_id", session.HallID,
			"movie_id", session.MovieID,
			"err", err,
		)
		return nil, err
	}

	return session, nil
}

func (s *sessionService) List() ([]models.Session, error) {

	sessions, err := s.sessionRepo.List()
	if err != nil {
		s.logger.Error(
			"failed to list sessions",
			"err", err,
		)
		return nil, err
	}

	return sessions, nil
}

func (s *sessionService) GetById(id uint) (*models.Session, error) {

	session, err := s.sessionRepo.GetById(id)
	if err != nil {
		s.logger.Warn(
			"session not found",
			"session_id", id,
			"error", err,
		)
		return nil, err
	}

	return session, nil
}

func (s *sessionService) Delete(id uint) error {

	if _, err := s.sessionRepo.GetById(id); err != nil {
		s.logger.Warn(
			"session not found",
			"session_id", id,
			"error", err,
		)
		return err
	}

	if err := s.sessionRepo.Delete(id); err != nil {
		s.logger.Error(
			"failed to delete session",
			"session_id", id,
			"err", err,
		)
		return err
	}
	s.logger.Info("session deleted successfully", "id", id)
	return nil
}

func (s *sessionService) ListByMovieID(movieID uint) ([]models.Session, error) {

	sessions, err := s.sessionRepo.ListByMovieID(movieID)
	if err != nil {
		s.logger.Error(
			"failed to list sessions by movie id",
			"movie_id", movieID,
			"err", err,
		)
		return nil, err
	}

	return sessions, nil
}

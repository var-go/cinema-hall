package repository

import (
	"cinema-service/internal/models"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type SessionRepository interface {
	Create(*models.Session) error
	List() ([]models.Session, error)
	Update(id uint, session *models.Session) error
	Delete(id uint) error
	GetById(id uint) (*models.Session, error)
	ListByMovieID(movieID uint) ([]models.Session, error)
}

type sessionRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewSessionRepository(db *gorm.DB, logger *slog.Logger) SessionRepository {
	return &sessionRepository{
		db:     db,
		logger: logger,
	}
}

func (r *sessionRepository) Create(session *models.Session) error {
	if session == nil {
		r.logger.Warn("attempt to create nil session")
		return errors.New("session is nil")
	}

	if err := r.db.Create(session).Error; err != nil {
		r.logger.Error("failed to create session", "err", err)
		return err
	}

	return nil
}

func (r *sessionRepository) List() ([]models.Session, error) {
	var sessions []models.Session

	if err := r.db.Find(&sessions).Error; err != nil {
		r.logger.Error("failed to fetch sessions", "err", err)
		return nil, err
	}

	return sessions, nil
}

func (r *sessionRepository) Update(id uint, session *models.Session) error {
	if session == nil {
		r.logger.Warn("attempt to update nil session")
		return errors.New("session is nil")
	}

	if err := r.db.
		Model(&models.Session{}).
		Where("id = ?", id).
		Updates(session).Error; err != nil {

		r.logger.Error(
			"failed to update session",
			"id", id,
			"err", err,
		)
		return err
	}

	return nil
}

func (r *sessionRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Session{}, id).Error; err != nil {
		r.logger.Error(
			"failed to delete session",
			"id", id,
			"err", err,
		)
		return err
	}

	return nil
}

func (r *sessionRepository) GetById(id uint) (*models.Session, error) {
	var session models.Session

	if err := r.db.First(&session, id).Error; err != nil {
		r.logger.Error(
			"failed to fetch session by id",
			"id", id,
			"err", err,
		)
		return nil, err
	}

	return &session, nil
}

func (r *sessionRepository) ListByMovieID(movieID uint) ([]models.Session, error) {
	var sessions []models.Session

	if err := r.db.
		Where("movie_id = ?", movieID).
		Find(&sessions).Error; err != nil {

		r.logger.Error(
			"failed to fetch sessions by movie id",
			"movie_id", movieID,
			"err", err,
		)
		return nil, err
	}

	return sessions, nil
}

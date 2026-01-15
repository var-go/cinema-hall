package repository

import (
	"cinema-service/internal/models"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type HallRepository interface {
	Create(*models.Hall) error
	List() ([]models.Hall, error)
	Update(id uint, hall *models.Hall) error
	Delete(id uint) error
	GetById(id uint) (*models.Hall, error)
}

type hallRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewHallRepository(db *gorm.DB, logger *slog.Logger) HallRepository {
	return &hallRepository{
		db:     db,
		logger: logger,
	}
}

func (r *hallRepository) Create(hall *models.Hall) error {
	if hall == nil {
		r.logger.Warn("attempt to create nil hall")
		return errors.New("hall is nil")
	}
	if err := r.db.Create(hall).Error; err != nil {
		r.logger.Error("failed to create a hall", "err", err)
		return err
	}
	return nil
}

func (r *hallRepository) List() ([]models.Hall, error) {
	var halls []models.Hall
	if err := r.db.Find(&halls).Error; err != nil {
		r.logger.Error("failed to fetch halls", "err", err)
		return nil, err
	}
	return halls, nil
}

func (r *hallRepository) Update(id uint, hall *models.Hall) error {

	if hall == nil {
		return errors.New("hall is nil")
	}
	return r.db.Model(&models.Hall{}).
		Where("id = ?", id).
		Updates(hall).Error
}

func (r *hallRepository) GetById(id uint) (*models.Hall, error) {
	var hall models.Hall

	if err := r.db.First(&hall, id).Error; err != nil {
		r.logger.Error("failed to fetch tool by id", "error", err, "id", id)
		return nil, err
	}
	return &hall, nil
}

func (r *hallRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Hall{}, id).Error; err != nil {
		r.logger.Error("failed to delete hall", "err", err)
		return err
	}
	return nil
}

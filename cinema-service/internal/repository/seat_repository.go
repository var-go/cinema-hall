package repository

import (
	"cinema-service/internal/models"
	"errors"
	"log/slog"

	"gorm.io/gorm"
)

type SeatRepository interface {
	Create(*models.Seat) error
	List() ([]models.Seat, error)
	Update(id uint, seat *models.Seat) error
	Delete(id uint) error
	GetById(id uint) (*models.Seat, error)
}

type seatRepository struct {
	db     *gorm.DB
	logger *slog.Logger
}

func NewSeatRepository(db *gorm.DB, logger *slog.Logger) SeatRepository {
	return &seatRepository{
		db:     db,
		logger: logger,
	}
}

func (r *seatRepository) Create(seat *models.Seat) error {
	if seat == nil {
		r.logger.Warn("attempt to create nil seat")
		return errors.New("seat is nil")
	}
	if err := r.db.Create(seat).Error; err != nil {
		r.logger.Error("failed to create a seat", "err", err)
		return err
	}
	return nil
}

func (r *seatRepository) List() ([]models.Seat, error) {
	var seats []models.Seat
	if err := r.db.Find(&seats).Error; err != nil {
		r.logger.Error("failed to fetch seats", "err", err)
		return nil, err
	}
	return seats, nil
}

func (r *seatRepository) Update(id uint, seat *models.Seat) error {
	if seat == nil {
		return errors.New("seat is nil")
	}

	return r.db.Model(&models.Seat{}).
		Where("id = ?", id).
		Updates(seat).Error
}

func (r *seatRepository) GetById(id uint) (*models.Seat, error) {
	var seat models.Seat

	if err := r.db.First(&seat, id).Error; err != nil {
		r.logger.Error("failed to fetch seat by id", "error", err, "id", id)
		return nil, err
	}
	return &seat, nil
}

func (r *seatRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Seat{}, id).Error; err != nil {
		r.logger.Error("failed to delete seat", "err", err)
		return err
	}
	return nil
}

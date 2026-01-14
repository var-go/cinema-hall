package repository

import (
	"booking-service/internal/models"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(booking *models.Booking) (*models.Booking, error)
	List() ([]models.Booking, error)
	GetByID(id uint) (*models.Booking, error)
	Update(id uint, req models.Booking) error
	Delete(id uint) error
}

type gormBookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &gormBookingRepository{
		db: db,
	}
}

func (r *gormBookingRepository) Create(booking *models.Booking) (*models.Booking, error) {

	if err := r.db.Create(&booking).Error; err != nil {
		log.Errorf("failed to create booking: %d", err)
	}

	return booking, nil
}

func (r *gormBookingRepository) List() ([]models.Booking, error) {
	var bookings []models.Booking

	if err := r.db.Find(&bookings).Error; err != nil {
		log.Errorf("failed to get bookings list")
		return nil, err
	}

	return bookings, nil
}

func (r *gormBookingRepository) GetByID(id uint) (*models.Booking, error) {
	var booking models.Booking

	if err := r.db.First(&booking, id).Error; err != nil {
		log.Errorf("failed to get booking by id")
		return nil, err
	}

	return &booking, nil
}

func (r *gormBookingRepository) Update(id uint, req models.Booking) error {
	if err := r.db.Model(&models.Booking{}).Where("id = ?", id).Updates(req).Error; err != nil {
		log.Error("error to update")
		return err
	}

	return nil
}

func (r *gormBookingRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Booking{}, id).Error; err != nil {
		log.Error("failed to remove booking")
		return err
	}

	return nil
}

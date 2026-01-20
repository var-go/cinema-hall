package repository

import (
	"booking-service/internal/constants"
	"booking-service/internal/models"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(tx *gorm.DB, booking *models.Booking) (*models.Booking, error)
	List() ([]models.Booking, error)
	GetByID(id uint) (*models.Booking, error)
	GetByIDWithTx(tx *gorm.DB, id uint) (*models.Booking, error)
	Update(id uint, req models.Booking) error
	UpdateWithTx(tx *gorm.DB, id uint, req models.Booking) error
	Delete(id uint) error
	CheckBooked(sessionID uint, seatsID []uint) ([]uint, error)
}

type gormBookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &gormBookingRepository{
		db: db,
	}
}

func (r *gormBookingRepository) Create(tx *gorm.DB, booking *models.Booking) (*models.Booking, error) {

	if err := tx.Create(&booking).Error; err != nil {
		log.Errorf("failed to create booking: %d", err)
		return nil, err
	}

	return booking, nil
}

func (r *gormBookingRepository) List() ([]models.Booking, error) {
	var bookings []models.Booking

	if err := r.db.Preload("BookedSeats").Find(&bookings).Error; err != nil {
		log.Errorf("failed to get bookings list")
		return nil, err
	}

	return bookings, nil
}

func (r *gormBookingRepository) GetByID(id uint) (*models.Booking, error) {
	var booking models.Booking

	if err := r.db.Preload("BookedSeats").First(&booking, id).Error; err != nil {
		log.Errorf("failed to get booking by id")
		return nil, err
	}

	return &booking, nil
}

func (r *gormBookingRepository) GetByIDWithTx(tx *gorm.DB, id uint) (*models.Booking, error) {
	var booking models.Booking

	if err := tx.Preload("BookedSeats").First(&booking, id).Error; err != nil {
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

func (r *gormBookingRepository) UpdateWithTx(tx *gorm.DB, id uint, req models.Booking) error {
	if err := tx.Model(&models.Booking{}).Where("id = ?", id).Updates(req).Error; err != nil {
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

func (r *gormBookingRepository) CheckBooked(sessionID uint, seatIDs []uint) ([]uint, error) {
	if len(seatIDs) == 0 {
		return []uint{}, nil
	}

	var bookedSeats []models.BookedSeat

	err := r.db.
		Joins("JOIN bookings ON booked_seats.booking_id = bookings.id").
		Where("bookings.session_id = ? AND bookings.booking_status IN (?, ?)",
			sessionID, constants.Pending, constants.Confirmed).
		Find(&bookedSeats).Error

	if err != nil {
		log.Errorf("failed to check booked seats: %v", err)
		return nil, err
	}

	var bookedSeatIDs = []uint{}

	for _, seat := range bookedSeats {
		bookedSeatIDs = append(bookedSeatIDs, seat.SeatID)
	}

	return bookedSeatIDs, nil
}

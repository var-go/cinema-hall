package repository

import (
	"booking-service/internal/config"
	"booking-service/internal/constants"
	"booking-service/internal/models"
	"time"

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
	FindExpiredPendingBookings() ([]models.Booking, error)
	FindBookingsForEndedSessions() ([]models.Booking, error)
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
		config.GetLogger().Error("Failed to create booking", "error", err, "session_id", booking.SessionID, "user_id", booking.UserID)
		return nil, err
	}

	return booking, nil
}

func (r *gormBookingRepository) List() ([]models.Booking, error) {
	var bookings []models.Booking

	if err := r.db.Preload("BookedSeats").Find(&bookings).Error; err != nil {
		config.GetLogger().Error("Failed to get bookings list", "error", err)
		return nil, err
	}

	return bookings, nil
}

func (r *gormBookingRepository) GetByID(id uint) (*models.Booking, error) {
	var booking models.Booking

	if err := r.db.Preload("BookedSeats").First(&booking, id).Error; err != nil {
		config.GetLogger().Error("Failed to get booking by id", "error", err, "booking_id", id)
		return nil, err
	}

	return &booking, nil
}

func (r *gormBookingRepository) GetByIDWithTx(tx *gorm.DB, id uint) (*models.Booking, error) {
	var booking models.Booking

	if err := tx.Preload("BookedSeats").First(&booking, id).Error; err != nil {
		config.GetLogger().Error("Failed to get booking by id in transaction", "error", err, "booking_id", id)
		return nil, err
	}

	return &booking, nil
}

func (r *gormBookingRepository) Update(id uint, req models.Booking) error {
	if err := r.db.Model(&models.Booking{}).Where("id = ?", id).Updates(req).Error; err != nil {
		config.GetLogger().Error("Failed to update booking", "error", err, "booking_id", id)
		return err
	}

	return nil
}

func (r *gormBookingRepository) UpdateWithTx(tx *gorm.DB, id uint, req models.Booking) error {
	if err := tx.Model(&models.Booking{}).Where("id = ?", id).Updates(req).Error; err != nil {
		config.GetLogger().Error("Failed to update booking in transaction", "error", err, "booking_id", id)
		return err
	}
	return nil
}

func (r *gormBookingRepository) Delete(id uint) error {
	if err := r.db.Delete(&models.Booking{}, id).Error; err != nil {
		config.GetLogger().Error("Failed to delete booking", "error", err, "booking_id", id)
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
		Where("bookings.session_id = ? AND bookings.booking_status IN (?, ?) AND booked_seats.seat_id IN ?",
			sessionID, constants.Pending, constants.Confirmed, seatIDs).
		Find(&bookedSeats).Error

	if err != nil {
		config.GetLogger().Error("Failed to check booked seats", "error", err, "session_id", sessionID, "seat_ids", seatIDs)
		return nil, err
	}

	var bookedSeatIDs = []uint{}

	for _, seat := range bookedSeats {
		bookedSeatIDs = append(bookedSeatIDs, seat.SeatID)
	}

	return bookedSeatIDs, nil
}

func (r *gormBookingRepository) FindExpiredPendingBookings() ([]models.Booking, error) {
	var bookings []models.Booking

	err := r.db.
		Where("booking_status = ? AND expires_at < ?", constants.Pending, time.Now()).
		Find(&bookings).Error

	if err != nil {
		config.GetLogger().Error("Failed to find expired pending bookings", "error", err)
		return nil, err
	}

	return bookings, nil
}

func (r *gormBookingRepository) FindBookingsForEndedSessions() ([]models.Booking, error) {
	var bookings []models.Booking

	err := r.db.
		Where("session_end_time < ? AND booking_status IN (?, ?)",
			time.Now(), constants.Pending, constants.Confirmed).
		Find(&bookings).Error

	if err != nil {
		config.GetLogger().Error("Failed to find bookings for ended sessions", "error", err)
		return nil, err
	}

	return bookings, nil
}

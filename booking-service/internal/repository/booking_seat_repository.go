package repository

import (
	"booking-service/internal/models"

	"github.com/gofiber/fiber/v2/log"
	"gorm.io/gorm"
)

type BookingSeatRepository interface {
	Create(bookingID uint, seatList []uint) error
}

type gormBookingSeat struct {
	db *gorm.DB
}

func NewBookingSeatRepository(db *gorm.DB) BookingSeatRepository {
	return &gormBookingSeat{
		db: db,
	}
}

func (r *gormBookingSeat) Create(bookingID uint, seatList []uint) error {
	var bookedSeats = make([]models.BookedSeat, 0, len(seatList))

	for _, seat := range seatList {
		bookedSeats = append(bookedSeats, models.BookedSeat{
			BookingID: bookingID,
			SeatID:    seat,
		})
	}

	if err := r.db.Create(&bookedSeats).Error; err != nil {
		log.Errorf("failed to create booked seats: %v", err)
		return err
	}

	return nil
}

package services

import (
	"booking-service/internal/clients"
	"booking-service/internal/config"
	"booking-service/internal/constants"
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"booking-service/internal/repository"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

type BookingService interface {
	Create(req dto.BookingCreateRequest) (*models.Booking, error)
	List() ([]models.Booking, error)
	GetByID(id uint) (*models.Booking, error)
	Update(id uint, req dto.BookingUpdateRequest) (*models.Booking, error)
	Delete(id uint) error

	ConfirmBooking(id uint) (*models.Booking, error)
	CancelBooking(id uint) (*models.Booking, error)
	ExpireOldBookings() error
	FreeSeatsForEndedSessions() error
	ExpireBooking(id uint) (*models.Booking, error)
}

type bookingService struct {
	bookingRepo     repository.BookingRepository
	bookingSeatRepo repository.BookingSeatRepository
	db              *gorm.DB
}

func NewBookingService(bookingRepo repository.BookingRepository, bookingSeatRepo repository.BookingSeatRepository, db *gorm.DB) BookingService {
	return &bookingService{
		bookingRepo:     bookingRepo,
		bookingSeatRepo: bookingSeatRepo,
		db:              db,
	}
}

func (s *bookingService) Create(req dto.BookingCreateRequest) (*models.Booking, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	session, err := clients.GetSession(req.SessionID)
	if err != nil {
		tx.Rollback()
		config.GetLogger().Error("Failed to get session", "error", err, "session_id", req.SessionID)
		return nil, fmt.Errorf("session not found")
	}

	if !session.StartTime.After(time.Now()) {
		tx.Rollback()
		return nil, fmt.Errorf("session already started")
	}

	bookedSeats, err := s.bookingRepo.CheckBooked(req.SessionID, req.SeatsID)
	if err != nil {
		tx.Rollback()
		config.GetLogger().Error("Failed to check booked seats", "error", err, "session_id", req.SessionID, "seats", req.SeatsID)
		return nil, err
	}
	if len(bookedSeats) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("seats already booked: %v", bookedSeats)
	}

	var booking = models.Booking{
		SessionID:        req.SessionID,
		UserID:           req.UserID,
		BookingStatus:    constants.Pending,
		PaymentStatus:    constants.PaymentPending,
		ExpiresAt:        time.Now().Add(constants.BookingTimeoutMinutes * time.Minute),
		SessionStartTime: session.StartTime,
		SessionEndTime:   session.EndTime,
	}

	newBooking, err := s.bookingRepo.Create(tx, &booking)
	if err != nil {
		tx.Rollback()
		config.GetLogger().Error("Failed to create booking", "error", err, "session_id", req.SessionID, "user_id", req.UserID)
		return nil, err
	}

	err = s.bookingSeatRepo.Create(tx, newBooking.ID, req.SeatsID)
	if err != nil {
		tx.Rollback()
		config.GetLogger().Error("Failed to create booked seats", "error", err, "booking_id", newBooking.ID, "seats", req.SeatsID)
		return nil, err
	}

	bookingWithSeats, err := s.bookingRepo.GetByIDWithTx(tx, newBooking.ID)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return bookingWithSeats, nil
}

func (s *bookingService) List() ([]models.Booking, error) {
	list, err := s.bookingRepo.List()
	if err != nil {
		return nil, err
	}

	return list, nil
}

func (s *bookingService) GetByID(id uint) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		config.GetLogger().Error("Failed to get booking by id", "error", err, "booking_id", id)
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) Update(id uint, req dto.BookingUpdateRequest) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		config.GetLogger().Error("Failed to get booking for update", "error", err, "booking_id", id)
		return nil, err
	}

	if req.BookingStatus != nil {
		booking.BookingStatus = *req.BookingStatus
	}

	if err := s.bookingRepo.Update(id, *booking); err != nil {
		config.GetLogger().Error("Failed to update booking", "error", err, "booking_id", id)
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) Delete(id uint) error {
	if err := s.bookingRepo.Delete(id); err != nil {
		config.GetLogger().Error("Failed to delete booking", "error", err, "booking_id", id)
		return err
	}

	return nil
}

func (s *bookingService) ConfirmBooking(id uint) (*models.Booking, error) {

	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			return nil, constants.ErrBookingNotFound
		}
		return nil, err
	}

	if !booking.ExpiresAt.After(time.Now()) {
		return nil, constants.ErrBookingExpired
	}

	switch booking.BookingStatus {
	case constants.Expired:
		return nil, constants.ErrBookingExpired
	case constants.Cancelled:
		return nil, constants.ErrBookingAlreadyCancelled
	case constants.Confirmed:
		return nil, constants.ErrBookingAlreadyConfirmed
	case constants.Pending:
		booking.BookingStatus = constants.Confirmed
		booking.PaymentStatus = constants.PaymentPaid
		err = s.bookingRepo.Update(booking.ID, *booking)
		if err != nil {
			return nil, err
		}
	default:
		return nil, constants.ErrInvalidBookingStatus
	}

	return booking, nil
}

func (s *bookingService) CancelBooking(id uint) (*models.Booking, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	booking, err := s.bookingRepo.GetByIDWithTx(tx, id)
	if err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			tx.Rollback()
			return nil, constants.ErrBookingNotFound
		}
		tx.Rollback()
		return nil, err
	}

	switch booking.BookingStatus {
	case constants.Expired:
		tx.Rollback()
		return nil, constants.ErrBookingExpired
	case constants.Cancelled:
		tx.Rollback()
		return nil, constants.ErrBookingAlreadyCancelled
	case constants.Pending:
		booking.BookingStatus = constants.Cancelled
	default:
		tx.Rollback()
		return nil, constants.ErrInvalidBookingStatus
	}

	err = s.bookingRepo.UpdateWithTx(tx, booking.ID, *booking)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.bookingSeatRepo.DeleteByBookingID(tx, booking.ID); err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}

func (s *bookingService) ExpireBooking(id uint) (*models.Booking, error) {
	tx := s.db.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r)
		}
	}()

	booking, err := s.bookingRepo.GetByIDWithTx(tx, id)
	if err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			tx.Rollback()
			return nil, constants.ErrBookingNotFound
		}
		tx.Rollback()
		return nil, err
	}

	if booking.BookingStatus != constants.Pending {
		tx.Rollback()
		return nil, fmt.Errorf("booking is not in pending status: %s", booking.BookingStatus)
	}

	booking.BookingStatus = constants.Expired

	err = s.bookingRepo.UpdateWithTx(tx, booking.ID, *booking)
	if err != nil {
		tx.Rollback()
		return nil, err
	}

	if err := s.bookingSeatRepo.DeleteByBookingID(tx, booking.ID); err != nil {
		tx.Rollback()
		config.GetLogger().Error("Failed to delete seats for expired booking",
			"error", err, "booking_id", booking.ID)
		return nil, err
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return booking, nil
}

func (s *bookingService) ExpireOldBookings() error {
	expiredBookings, err := s.bookingRepo.FindExpiredPendingBookings()
	if err != nil {
		config.GetLogger().Error("Failed to find expired bookings", "error", err)
		return err
	}

	if len(expiredBookings) == 0 {
		return nil
	}

	config.GetLogger().Info("Found expired bookings to expire", "count", len(expiredBookings))

	for _, booking := range expiredBookings {
		_, err := s.ExpireBooking(booking.ID)
		if err != nil {
			config.GetLogger().Error("Failed to expire booking",
				"error", err, "booking_id", booking.ID)
			continue
		}

		config.GetLogger().Info("Expired booking processed and seats freed",
			"booking_id", booking.ID, "session_id", booking.SessionID)
	}

	return nil
}

func (s *bookingService) FreeSeatsForEndedSessions() error {
	endedSessionsBookings, err := s.bookingRepo.FindBookingsForEndedSessions()
	if err != nil {
		config.GetLogger().Error("Failed to find bookings for ended sessions", "error", err)
		return err
	}

	if len(endedSessionsBookings) == 0 {
		return nil
	}

	config.GetLogger().Info("Found bookings for ended sessions to free seats", "count", len(endedSessionsBookings))

	for _, booking := range endedSessionsBookings {
		tx := s.db.Begin()
		if tx.Error != nil {
			config.GetLogger().Error("Failed to start transaction", "error", tx.Error, "booking_id", booking.ID)
			continue
		}

		currentBooking, err := s.bookingRepo.GetByIDWithTx(tx, booking.ID)
		if err != nil {
			tx.Rollback()
			config.GetLogger().Error("Failed to get booking for ended session",
				"error", err, "booking_id", booking.ID)
			continue
		}

		var finalStatus constants.BookingStatus
		switch currentBooking.BookingStatus {
		case constants.Pending:
			finalStatus = constants.Expired
		case constants.Confirmed:
			finalStatus = constants.Finished
		default:
			tx.Rollback()
			config.GetLogger().Info("Booking already processed, skipping",
				"booking_id", booking.ID, "status", currentBooking.BookingStatus)
			continue
		}

		currentBooking.BookingStatus = finalStatus
		if err := s.bookingRepo.UpdateWithTx(tx, currentBooking.ID, *currentBooking); err != nil {
			tx.Rollback()
			config.GetLogger().Error("Failed to update booking status for ended session",
				"error", err, "booking_id", booking.ID, "final_status", finalStatus)
			continue
		}

		if err := s.bookingSeatRepo.DeleteByBookingID(tx, currentBooking.ID); err != nil {
			tx.Rollback()
			config.GetLogger().Error("Failed to delete seats for ended session",
				"error", err, "booking_id", booking.ID, "session_id", currentBooking.SessionID)
			continue
		}

		if err := tx.Commit().Error; err != nil {
			tx.Rollback()
			config.GetLogger().Error("Failed to commit transaction for ended session",
				"error", err, "booking_id", booking.ID)
			continue
		}

		config.GetLogger().Info("Seats freed for ended session",
			"booking_id", booking.ID,
			"session_id", currentBooking.SessionID,
			"old_status", booking.BookingStatus,
			"new_status", finalStatus)
	}

	return nil
}

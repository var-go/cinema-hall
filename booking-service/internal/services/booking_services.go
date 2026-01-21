package services

import (
	"booking-service/internal/clients"
	"booking-service/internal/constants"
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"booking-service/internal/repository"
	"errors"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2/log"
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
		log.Errorf("failed to get session: %v", err)
		return nil, fmt.Errorf("session not found")
	}

	if !session.StartTime.After(time.Now()) {
		tx.Rollback()
		return nil, fmt.Errorf("session already started")
	}

	bookedSeats, err := s.bookingRepo.CheckBooked(req.SessionID, req.SeatsID)
	if err != nil {
		tx.Rollback()
		log.Errorf("failed to check booked seats: %v", err)
		return nil, err
	}
	if len(bookedSeats) > 0 {
		tx.Rollback()
		return nil, fmt.Errorf("seats already booked: %v", bookedSeats)
	}

	var booking = models.Booking{
		SessionID:     req.SessionID,
		UserID:        req.UserID,
		BookingStatus: constants.Pending,
		PaymentStatus: constants.PaymentPending,
		ExpiresAt:     time.Now().Add(constants.BookingTimeoutMinutes * time.Minute),
	}

	newBooking, err := s.bookingRepo.Create(tx, &booking)
	if err != nil {
		tx.Rollback()
		log.Errorf("failed to create: %d", err)
		return nil, err
	}

	err = s.bookingSeatRepo.Create(tx, newBooking.ID, req.SeatsID)
	if err != nil {
		tx.Rollback()
		log.Errorf("failed to create booked seats: %v", err)
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
		log.Error("failed to get by id")
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) Update(id uint, req dto.BookingUpdateRequest) (*models.Booking, error) {
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		log.Error("record not found")
		return nil, err
	}

	if req.BookingStatus != nil {
		booking.BookingStatus = *req.BookingStatus
	}

	if err := s.bookingRepo.Update(id, *booking); err != nil {
		log.Error("failed to update booking")
		return nil, err
	}

	return booking, nil
}

func (s *bookingService) Delete(id uint) error {
	if err := s.bookingRepo.Delete(id); err != nil {
		log.Error("failed to remove booking")
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

package services

import (
	"booking-service/internal/constants"
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"booking-service/internal/repository"
	"errors"

	"github.com/gofiber/fiber/v2/log"
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
	bookingRepo repository.BookingRepository
}

func NewBookingService(bookingRepo repository.BookingRepository) BookingService {
	return &bookingService{
		bookingRepo: bookingRepo,
	}
}

func (s *bookingService) Create(req dto.BookingCreateRequest) (*models.Booking, error) {
	var booking = models.Booking{
		SessionID: req.SessionID,
		UserID:    req.UserID,
	}

	newBooking, err := s.bookingRepo.Create(&booking)
	if err != nil {
		log.Errorf("failed to create: %d", err)
		return nil, err
	}

	return newBooking, nil
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

	switch booking.BookingStatus {
	case constants.Expired:
		return nil, constants.ErrBookingExpired
	case constants.Cancelled:
		return nil, constants.ErrBookingAlreadyCancelled
	case constants.Confirmed:
		return nil, constants.ErrBookingAlreadyConfirmed
	case constants.Pending:
		booking.BookingStatus = constants.Confirmed
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
	booking, err := s.bookingRepo.GetByID(id)
	if err != nil {
		if errors.Is(err, constants.ErrBookingNotFound) {
			return nil, constants.ErrBookingNotFound
		}
		return nil, err
	}

	switch booking.BookingStatus {
	case constants.Expired:
		return nil, constants.ErrBookingExpired
	case constants.Cancelled:
		return nil, constants.ErrBookingAlreadyCancelled
	case constants.Pending:
		booking.BookingStatus = constants.Cancelled
	default:
		return nil, constants.ErrInvalidBookingStatus
	}

	err = s.bookingRepo.Update(booking.ID, *booking)
	if err != nil {
		return nil, err
	}

	return booking, nil
}

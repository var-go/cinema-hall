package services

import (
	"booking-service/internal/dto"
	"booking-service/internal/models"
	"booking-service/internal/repository"

	"github.com/gofiber/fiber/v2/log"
)

type BookingService interface {
	Create(req dto.BookingCreateRequest) (*models.Booking, error)
	List() ([]models.Booking, error)
	GetByID(id uint) (*models.Booking, error)
	Update(id uint, req dto.BookingUpdateRequest) (*models.Booking, error)
	Delete(id uint) error
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
		CinemaID: req.CinemaID,
		UserID:   req.UserID,
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

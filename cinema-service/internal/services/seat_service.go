package services

import (
	"cinema-service/internal/dto"
	"cinema-service/internal/models"
	"cinema-service/internal/repository"
	"log/slog"
)

type SeatService interface {
	Create(hallID uint, req dto.CreateSeatRequest) (*models.Seat, error)
	UpdateSeat(id uint, req dto.UpdateSeatRequest) (*models.Seat, error)
	List() ([]models.Seat, error)
	Delete(id uint) error
}

type seatService struct {
	seatRepo repository.SeatRepository
	hallRepo repository.HallRepository
	logger   *slog.Logger
}

func NewSeatService(
	seatRepo repository.SeatRepository,
	hallRepo repository.HallRepository,
	logger *slog.Logger,
) SeatService {
	return &seatService{
		seatRepo: seatRepo,
		hallRepo: hallRepo,
		logger:   logger,
	}
}

func (s *seatService) Create(hallID uint, req dto.CreateSeatRequest) (*models.Seat, error) {

	if _, err := s.hallRepo.GetById(hallID); err != nil {
		s.logger.Warn(
			"hall not found while creating seat",
			"hall_id", hallID,
			"error", err,
		)
		return nil, err
	}

	seat := &models.Seat{
		HallID: hallID,
		Row:    req.Row,
		Number: req.Number,
		Type:   req.Type,
	}
	if err := s.seatRepo.Create(seat); err != nil {
		s.logger.Error(
			"failed to create seat", "err", err)
		return nil, err
	}

	return seat, nil
}

func (s *seatService) UpdateSeat(id uint, req dto.UpdateSeatRequest) (*models.Seat, error) {
	seat, err := s.seatRepo.GetById(id)
	if err != nil {
		s.logger.Warn("seat not found", "seat_id", id)
		return nil, err
	}

	if req.Row != nil {
		seat.Row = *req.Row
	}
	if req.Number != nil {
		seat.Number = *req.Number
	}
	if req.Type != nil {
		seat.Type = *req.Type
	}

	if err := s.seatRepo.Update(id, seat); err != nil {
		s.logger.Error(
			"failed to update seat",
			"seat_id", id,
			"hall_id", seat.HallID,
			"row", seat.Row,
			"number", seat.Number,
			"err", err,
		)
		return nil, err
	}

	return seat, nil
}

func (s *seatService) List() ([]models.Seat, error) {
	seats, err := s.seatRepo.List()
	if err != nil {
		s.logger.Error("service: failed to list seats", "err", err)
		return nil, err
	}
	if len(seats) == 0 {
		s.logger.Info("service: no seats")
	}
	return seats, nil
}

func (s seatService) Delete(id uint) error {
	err := s.seatRepo.Delete(id)
	if err != nil {
		s.logger.Error("failed to delete seat", "id", id)
		return err
	}
	s.logger.Info("seat deleted successfully", "id", id)
	return nil
}

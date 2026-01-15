package services

import (
	"cinema-service/internal/dto"
	"cinema-service/internal/models"
	"cinema-service/internal/repository"
	"log/slog"

	"gorm.io/gorm"
)

type HallService interface {
	CreateHall(req dto.CreateHallRequest) (*models.Hall, error)
	ListHall() ([]models.Hall, error)
	UpdateHall(id uint, req dto.UpdateHallRequest) (*models.Hall, error)
	GetHallByID(id uint) (*models.Hall, error)
	DeleteHall(id uint) error
}

type hallService struct {
	hallRepo repository.HallRepository
	logger   *slog.Logger
}

func NewHallService(
	hallRepo repository.HallRepository,
	logger *slog.Logger,
) HallService {
	return &hallService{
		hallRepo: hallRepo,
		logger:   logger,
	}
}

func (s *hallService) UpdateHall(id uint, req dto.UpdateHallRequest) (*models.Hall, error) {
	hall, err := s.hallRepo.GetById(id)
	if err != nil {
		return nil, err
	}

	if hall == nil {
		return nil, gorm.ErrRecordNotFound
	}

	if req.Number != nil {
		hall.Number = *req.Number
	}

	if err := s.hallRepo.Update(id, hall); err != nil {
		return nil, err
	}
	return hall, nil
}

func (s *hallService) CreateHall(req dto.CreateHallRequest) (*models.Hall, error) {
	hall := models.Hall{
		Number: req.Number,
	}
	if err := s.hallRepo.Create(&hall); err != nil {
		s.logger.Error("service: failed to create hall", "err", err)
		return nil, err
	}
	return &hall, nil
}

func (s *hallService) ListHall() ([]models.Hall, error) {
	halls, err := s.hallRepo.List()
	if err != nil {
		s.logger.Error("service: failed to list halls", "err", err)
		return nil, err
	}
	if len(halls) == 0 {
		s.logger.Info("service: no halls")
	}
	return halls, nil
}

func (s *hallService) GetHallByID(id uint) (*models.Hall, error) {
	hall, err := s.hallRepo.GetById(id)
	if err != nil {
		s.logger.Error("service: failed to fetch hall by ID", "err", err)
		return nil, err
	}
	return hall, nil
}

func (s *hallService) DeleteHall(id uint) error {
	err := s.hallRepo.Delete(id)
	if err != nil {
		s.logger.Error("failed to delete hall", "id", id)
		return err
	}
	s.logger.Info("hall deleted successfully", "id", id)
	return nil
}

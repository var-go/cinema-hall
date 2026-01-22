package services

import (
	"log/slog"
	"user-service/internal/auth"
	"user-service/internal/dto"
	"user-service/internal/errors"
	"user-service/internal/kafka"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*models.User, error)
	Login(req dto.LoginRequest) (string, error)
}

type authService struct {
	repo     repository.UserRepository
	producer *kafka.Producer
	log      *slog.Logger
}

func NewAuthService(repo repository.UserRepository, producer *kafka.Producer, log *slog.Logger) AuthService {
	return &authService{repo: repo, producer: producer, log: log}
}

func (s *authService) Register(req dto.RegisterRequest) (*models.User, error) {
	if _, err := s.repo.GetByEmail(req.Email); err == nil {
		return nil, errors.ErrUserAlreadyExists
	}

	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password),
		bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     "user",
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	if err := s.producer.SendUserCreated(kafka.UserCreatedEvent{
		ID:    user.ID,
		Email: user.Email,
	}); err != nil {
		s.log.Error("failed to send user.created event", "user_id", user.ID, "err", err)
	}

	return user, nil
}

func (s *authService) Login(req dto.LoginRequest) (string, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.Password),
	); err != nil {
		return "", err
	}

	return auth.GenerateToken(user.ID, user.Role)
}

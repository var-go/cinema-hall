package services

import (
	"user-service/internal/dto"
	"user-service/internal/models"
	"user-service/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Create(req dto.CreateUserRequest) (*models.User, error)
	Get(id uint) (*models.User, error)
	List() ([]models.User, error)
	Update(id uint, req dto.UpdateUserRequest) (*models.User, error)
	Delete(id uint) error
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

func (s *userService) Create(req dto.CreateUserRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword(
		[]byte(req.Password), bcrypt.DefaultCost,
	)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Email:    req.Email,
		Password: string(hashedPassword),
		Name:     req.Name,
		Role:     req.Role,
	}

	if user.Role == "" {
		user.Role = "user"
	}

	return user, s.repo.Create(user)
}

func (s *userService) Get(id uint) (*models.User, error) {
	return s.repo.GetByID(id)
}

func (s *userService) List() ([]models.User, error) {
	return s.repo.GetAll()
}

func (s *userService) Update(id uint, req dto.UpdateUserRequest) (*models.User, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if req.Email != nil {
		user.Email = *req.Email
	}

	if req.Role != nil {
		user.Role = *req.Role
	}

	return user, s.repo.Update(user)
}

func (s *userService) Delete(id uint) error {
	_, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(id)
}

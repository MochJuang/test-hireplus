package service

import (
	"errors"
	"hireplus-project/internal/config"
	"hireplus-project/internal/entity"
	"hireplus-project/internal/repository"
	"hireplus-project/internal/utils"
	"time"

	"github.com/google/uuid"
	e "hireplus-project/internal/exception"
)

type UserService interface {
	Register(firstName, lastName, phone, address, pin string) (*entity.User, error)
	Login(phone, pin string) (string, string, error)
	UpdateProfile(userID, firstName, lastName, address string) (*entity.User, error)
	GetUserBalance(userID string) (float64, error)
}

type userService struct {
	userRepo repository.UserRepository
	config   config.Config
}

func NewUserService(userRepo repository.UserRepository, cfg config.Config) UserService {
	return &userService{userRepo, cfg}
}

func (s *userService) Register(firstName, lastName, phone, address, pin string) (*entity.User, error) {
	user := &entity.User{
		ID:          uuid.New().String(),
		FirstName:   firstName,
		LastName:    lastName,
		PhoneNumber: phone,
		Address:     address,
		Pin:         pin,
		CreatedAt:   time.Now(),
	}
	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, e.Internal(err)
	}

	if err := s.userRepo.CreateUserBalance(user.ID); err != nil {
		return nil, e.Internal(err)
	}

	return user, nil
}

func (s *userService) Login(phone, pin string) (string, string, error) {
	user, err := s.userRepo.GetUserByPhone(phone)
	if err != nil {
		return "", "", e.Validation(errors.New("phone number or PIN is incorrect"))
	}

	if user.Pin != pin {
		return "", "", e.Validation(errors.New("phone number or PIN is incorrect"))
	}

	jwtKey := s.config.JWTSecret

	accessToken, err := utils.GenerateToken(user.ID, jwtKey)
	if err != nil {
		return "", "", e.Internal(err)
	}

	refreshToken, err := utils.GenerateToken(user.ID, jwtKey)
	if err != nil {
		return "", "", e.Internal(err)
	}

	return accessToken, refreshToken, nil
}

func (s *userService) UpdateProfile(userID, firstName, lastName, address string) (*entity.User, error) {
	user, err := s.userRepo.GetUserByID(userID)
	if err != nil {
		return nil, e.Internal(err)
	}

	user.FirstName = firstName
	user.LastName = lastName
	user.Address = address
	user.UpdatedAt = time.Now()

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, e.Internal(err)
	}

	return user, nil
}

func (s *userService) GetUserBalance(userID string) (float64, error) {
	return s.userRepo.GetUserBalance(userID)
}

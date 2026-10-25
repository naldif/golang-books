package auth

import (
	"errors"

	"golang-buku/internal/modules/users"
	apperror "golang-buku/pkg/error"
	"golang-buku/utils"

	"gorm.io/gorm"
)

type Service interface {
	Register(input RegisterInput) (*users.User, error)
	Login(input LoginInput) (string, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Register(input RegisterInput) (*users.User, error) {
	existing, err := s.repo.FindByUsername(input.Username)
	if err == nil && existing != nil {
		return nil, apperror.ErrUserExists
	}
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		return nil, err
	}

	user := users.User{
		Username:  input.Username,
		Password:  hashedPassword,
		CreatedBy: input.Username,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *service) Login(input LoginInput) (string, error) {
	user, err := s.repo.FindByUsername(input.Username)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", apperror.ErrInvalidLogin
		}
		return "", err
	}

	if !utils.CheckPassword(input.Password, user.Password) {
		return "", apperror.ErrInvalidLogin
	}

	token, err := utils.GenerateToken(user.ID, user.Username)
	if err != nil {
		return "", err
	}

	return token, nil
}

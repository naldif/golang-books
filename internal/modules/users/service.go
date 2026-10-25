package users

import (
	"errors"

	apperror "golang-buku/pkg/error"
	"golang-buku/utils"

	"gorm.io/gorm"
)

type Service interface {
	Create(input CreateUserInput, createdBy string) (*UserResponse, error)
	FindAll() ([]UserResponse, error)
	FindByID(id uint) (*UserResponse, error)
	Update(id uint, input UpdateUserInput, modifiedBy string) (*UserResponse, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(input CreateUserInput, createdBy string) (*UserResponse, error) {
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

	user := User{
		Username:  input.Username,
		Password:  hashedPassword,
		CreatedBy: createdBy,
	}

	if err := s.repo.Create(&user); err != nil {
		return nil, err
	}

	res := formatUserResponse(&user)
	return &res, nil
}

func (s *service) FindAll() ([]UserResponse, error) {
	userList, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	responses := make([]UserResponse, 0, len(userList))
	for _, u := range userList {
		responses = append(responses, formatUserResponse(&u))
	}
	return responses, nil
}

func (s *service) FindByID(id uint) (*UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	res := formatUserResponse(user)
	return &res, nil
}

func (s *service) Update(id uint, input UpdateUserInput, modifiedBy string) (*UserResponse, error) {
	user, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrUserNotFound
		}
		return nil, err
	}

	if input.Username != "" && input.Username != user.Username {
		existing, err := s.repo.FindByUsername(input.Username)
		if err == nil && existing != nil && existing.ID != user.ID {
			return nil, apperror.ErrUserExists
		}
		user.Username = input.Username
	}

	if input.Password != "" {
		hashed, err := utils.HashPassword(input.Password)
		if err != nil {
			return nil, err
		}
		user.Password = hashed
	}

	user.ModifiedBy = modifiedBy

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	res := formatUserResponse(user)
	return &res, nil
}

func (s *service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrUserNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

func formatUserResponse(u *User) UserResponse {
	return UserResponse{
		ID:         u.ID,
		Username:   u.Username,
		CreatedAt:  u.CreatedAt,
		CreatedBy:  u.CreatedBy,
		ModifiedAt: u.ModifiedAt,
		ModifiedBy: u.ModifiedBy,
	}
}

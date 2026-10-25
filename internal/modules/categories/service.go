package categories

import (
	"errors"

	apperror "golang-buku/pkg/error"

	"gorm.io/gorm"
)

type Service interface {
	GetAll() ([]Category, error)
	GetByID(id uint) (*Category, error)
	Create(input CreateCategoryInput, username string) (*Category, error)
	Update(id uint, input UpdateCategoryInput, username string) (*Category, error)
	Delete(id uint) error
	GetBooksByCategoryID(id uint) ([]map[string]any, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAll() ([]Category, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id uint) (*Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrCategoryNotFound
		}
		return nil, err
	}
	return category, nil
}

func (s *service) Create(input CreateCategoryInput, username string) (*Category, error) {
	category := Category{
		Name:      input.Name,
		CreatedBy: username,
	}

	if err := s.repo.Create(&category); err != nil {
		return nil, err
	}
	return &category, nil
}

func (s *service) Update(id uint, input UpdateCategoryInput, username string) (*Category, error) {
	category, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrCategoryNotFound
		}
		return nil, err
	}

	category.Name = input.Name
	category.ModifiedBy = username

	if err := s.repo.Update(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (s *service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrCategoryNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

func (s *service) GetBooksByCategoryID(id uint) ([]map[string]any, error) {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrCategoryNotFound
		}
		return nil, err
	}

	return s.repo.FindBooksByCategoryID(id)
}

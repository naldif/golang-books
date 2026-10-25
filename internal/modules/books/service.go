package books

import (
	"errors"

	apperror "golang-buku/pkg/error"

	"gorm.io/gorm"
)

type Service interface {
	GetAll() ([]Book, error)
	GetByID(id uint) (*Book, error)
	Create(input CreateBookInput, username string) (*Book, error)
	Update(id uint, input UpdateBookInput, username string) (*Book, error)
	Delete(id uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func calculateThickness(totalPage int) string {
	if totalPage > 100 {
		return "tebal"
	}
	return "tipis"
}

func validateReleaseYear(year int) error {
	if year < 1980 || year > 2024 {
		return apperror.ErrInvalidYear
	}
	return nil
}

func (s *service) GetAll() ([]Book, error) {
	return s.repo.FindAll()
}

func (s *service) GetByID(id uint) (*Book, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrBookNotFound
		}
		return nil, err
	}
	return book, nil
}

func (s *service) Create(input CreateBookInput, username string) (*Book, error) {
	if err := validateReleaseYear(input.ReleaseYear); err != nil {
		return nil, err
	}

	thickness := calculateThickness(input.TotalPage)

	book := Book{
		Title:       input.Title,
		Description: input.Description,
		ImageURL:    input.ImageURL,
		ReleaseYear: input.ReleaseYear,
		Price:       input.Price,
		TotalPage:   input.TotalPage,
		Thickness:   thickness,
		CategoryID:  input.CategoryID,
		CreatedBy:   username,
	}

	if err := s.repo.Create(&book); err != nil {
		return nil, err
	}

	return s.repo.FindByID(book.ID)
}

func (s *service) Update(id uint, input UpdateBookInput, username string) (*Book, error) {
	book, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.ErrBookNotFound
		}
		return nil, err
	}

	if err := validateReleaseYear(input.ReleaseYear); err != nil {
		return nil, err
	}

	thickness := calculateThickness(input.TotalPage)

	book.Title = input.Title
	book.Description = input.Description
	book.ImageURL = input.ImageURL
	book.ReleaseYear = input.ReleaseYear
	book.Price = input.Price
	book.TotalPage = input.TotalPage
	book.Thickness = thickness
	book.CategoryID = input.CategoryID
	book.ModifiedBy = username

	if err := s.repo.Update(book); err != nil {
		return nil, err
	}

	return s.repo.FindByID(book.ID)
}

func (s *service) Delete(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return apperror.ErrBookNotFound
		}
		return err
	}

	return s.repo.Delete(id)
}

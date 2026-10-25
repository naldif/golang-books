package books

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Book, error)
	FindByID(id uint) (*Book, error)
	Create(book *Book) error
	Update(book *Book) error
	Delete(id uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Book, error) {
	var books []Book
	err := r.db.Preload("Category").Order("id desc").Find(&books).Error
	return books, err
}

func (r *repository) FindByID(id uint) (*Book, error) {
	var book Book
	err := r.db.Preload("Category").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *repository) Create(book *Book) error {
	return r.db.Create(book).Error
}

func (r *repository) Update(book *Book) error {
	return r.db.Save(book).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Book{}, id).Error
}

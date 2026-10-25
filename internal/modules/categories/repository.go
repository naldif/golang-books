package categories

import (
	"gorm.io/gorm"
)

type Repository interface {
	FindAll() ([]Category, error)
	FindByID(id uint) (*Category, error)
	Create(category *Category) error
	Update(category *Category) error
	Delete(id uint) error
	FindBooksByCategoryID(categoryID uint) ([]map[string]any, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindAll() ([]Category, error) {
	var categories []Category
	err := r.db.Order("id desc").Find(&categories).Error
	return categories, err
}

func (r *repository) FindByID(id uint) (*Category, error) {
	var category Category
	err := r.db.First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *repository) Create(category *Category) error {
	return r.db.Create(category).Error
}

func (r *repository) Update(category *Category) error {
	return r.db.Save(category).Error
}

func (r *repository) Delete(id uint) error {
	return r.db.Delete(&Category{}, id).Error
}

func (r *repository) FindBooksByCategoryID(categoryID uint) ([]map[string]any, error) {
	var books []map[string]any
	err := r.db.Table("books").Where("category_id = ?", categoryID).Find(&books).Error
	return books, err
}

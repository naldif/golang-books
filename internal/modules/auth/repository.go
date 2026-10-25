package auth

import (
	"golang-buku/internal/modules/users"

	"gorm.io/gorm"
)

type Repository interface {
	FindByUsername(username string) (*users.User, error)
	Create(user *users.User) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindByUsername(username string) (*users.User, error) {
	var user users.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *repository) Create(user *users.User) error {
	return r.db.Create(user).Error
}

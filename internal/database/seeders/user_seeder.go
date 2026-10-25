package seeders

import (
	"fmt"
	"log"

	"golang-buku/internal/modules/users"
	"golang-buku/utils"

	"gorm.io/gorm"
)

func SeedUser(db *gorm.DB) error {
	var count int64
	if err := db.Model(&users.User{}).Where("username = ?", "nito").Count(&count).Error; err != nil {
		return err
	}

	if count > 0 {
		return nil
	}

	hashedPassword, err := utils.HashPassword("password123!")
	if err != nil {
		return fmt.Errorf("gagal hashing password: %w", err)
	}

	defaultUser := users.User{
		Username:  "nito",
		Password:  hashedPassword,
		CreatedBy: "system_seeder",
	}

	if err := db.Create(&defaultUser).Error; err != nil {
		return fmt.Errorf("gagal menyimpan user seeder: %w", err)
	}

	log.Println("User seeder 'nito' berhasil dijalankan")
	return nil
}

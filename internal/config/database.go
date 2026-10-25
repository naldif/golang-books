package config

import (
	"log"
	"os"

	"golang-buku/internal/database/seeders"
	"golang-buku/internal/modules/books"
	"golang-buku/internal/modules/categories"
	"golang-buku/internal/modules/users"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() *gorm.DB {
	LoadEnv()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=db_book port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	err = DB.AutoMigrate(
		&users.User{},
		&categories.Category{},
		&books.Book{},
	)
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi tabel database: %v", err)
	}

	if err := seeders.SeedUser(DB); err != nil {
		log.Printf("Gagal menjalankan seeder user: %v\n", err)
	}

	return DB
}

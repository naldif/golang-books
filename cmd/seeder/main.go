package main

import (
	"log"

	"golang-buku/internal/config"
	"golang-buku/internal/database/seeders"
)

func main() {
	config.LoadEnv()

	db := config.InitDB()

	if err := seeders.SeedUser(db); err != nil {
		log.Fatalf("Gagal menjalankan seeder user: %v", err)
	}
}

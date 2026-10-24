package main

import (
	_ "embed"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

//go:embed openapi.json
var openapiSpec []byte

// @title API Siswa
// @version 1.0
// @description Dokumentasi REST API Siswa - Bootcamp Sanbercode Golang
// @host localhost:8081
// @BasePath /

type Siswa struct {
	ID    uint   `gorm:"primaryKey" json:"id" example:"1"`
	Nama  string `gorm:"not null" json:"nama" binding:"required" example:"Budi Santoso"`
	NIS   string `gorm:"not null" json:"nis" binding:"required" example:"12345"`
	Kelas string `gorm:"not null" json:"kelas" binding:"required" example:"XII-RPL-1"`
}

type CreateSiswaInput struct {
	Nama  string `json:"nama" binding:"required" example:"Budi Santoso"`
	NIS   string `json:"nis" binding:"required" example:"12345"`
	Kelas string `json:"kelas" binding:"required" example:"XII-RPL-1"`
}

type APIResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
	Error   string `json:"error,omitempty"`
}

var DB *gorm.DB

func initDB() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		dsn = "host=localhost user=postgres password=postgres dbname=siswa port=5432 sslmode=disable TimeZone=Asia/Jakarta"
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Gagal terhubung ke database: %v", err)
	}

	err = DB.AutoMigrate(&Siswa{})
	if err != nil {
		log.Fatalf("Gagal melakukan migrasi tabel: %v", err)
	}

	fmt.Println("Koneksi database dan migrasi berhasil!")
}

// @Summary Menambahkan siswa baru
// @Description Menyimpan data siswa baru ke database
// @Tags Siswa
// @Accept json
// @Produce json
// @Param siswa body CreateSiswaInput true "Data Siswa"
// @Success 201 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /siswa [post]
func createSiswa(c *gin.Context) {
	var input CreateSiswaInput

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "error",
			"message": "Validasi gagal: Nama, NIS, dan Kelas wajib diisi",
			"error":   err.Error(),
		})
		return
	}

	siswa := Siswa{
		Nama:  input.Nama,
		NIS:   input.NIS,
		Kelas: input.Kelas,
	}

	if result := DB.Create(&siswa); result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "error",
			"message": "Gagal menyimpan data siswa",
			"error":   result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, APIResponse{
		Status:  "success",
		Message: "Siswa berhasil ditambahkan",
		Data:    siswa,
	})
}

// @Summary Mengambil semua siswa
// @Description Mengambil semua daftar siswa dari database
// @Tags Siswa
// @Produce json
// @Success 200 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /siswa [get]
func getAllSiswa(c *gin.Context) {
	var daftarSiswa []Siswa

	if err := DB.Find(&daftarSiswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Gagal mengambil daftar siswa",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Daftar siswa berhasil diambil",
		Data:    daftarSiswa,
	})
}

// @Summary Detail siswa berdasarkan ID
// @Description Mengambil satu data siswa berdasarkan ID
// @Tags Siswa
// @Produce json
// @Param id path int true "ID Siswa"
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /siswa/{id} [get]
func getSiswaByID(c *gin.Context) {
	id := c.Param("id")
	var siswa Siswa

	if err := DB.First(&siswa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: fmt.Sprintf("Siswa dengan ID %s tidak ditemukan", id),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Terjadi kesalahan pada server",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Detail siswa berhasil ditemukan",
		Data:    siswa,
	})
}

// @Summary Memperbarui data siswa
// @Description Mengubah data siswa berdasarkan ID
// @Tags Siswa
// @Accept json
// @Produce json
// @Param id path int true "ID Siswa"
// @Param siswa body CreateSiswaInput true "Data Siswa yang diperbarui"
// @Success 200 {object} APIResponse
// @Failure 400 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /siswa/{id} [put]
func updateSiswa(c *gin.Context) {
	id := c.Param("id")
	var siswa Siswa

	if err := DB.First(&siswa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: fmt.Sprintf("Siswa dengan ID %s tidak ditemukan", id),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Terjadi kesalahan pada server",
			Error:   err.Error(),
		})
		return
	}

	var input CreateSiswaInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, APIResponse{
			Status:  "error",
			Message: "Validasi gagal: Nama, NIS, dan Kelas wajib diisi",
			Error:   err.Error(),
		})
		return
	}

	siswa.Nama = input.Nama
	siswa.NIS = input.NIS
	siswa.Kelas = input.Kelas

	if err := DB.Save(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Gagal memperbarui data siswa",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: "Data siswa berhasil diperbarui",
		Data:    siswa,
	})
}

// @Summary Menghapus siswa
// @Description Menghapus data siswa dari database berdasarkan ID
// @Tags Siswa
// @Produce json
// @Param id path int true "ID Siswa"
// @Success 200 {object} APIResponse
// @Failure 404 {object} APIResponse
// @Failure 500 {object} APIResponse
// @Router /siswa/{id} [delete]
func deleteSiswa(c *gin.Context) {
	id := c.Param("id")
	var siswa Siswa

	if err := DB.First(&siswa, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, APIResponse{
				Status:  "error",
				Message: fmt.Sprintf("Siswa dengan ID %s tidak ditemukan", id),
			})
			return
		}

		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Terjadi kesalahan pada server",
			Error:   err.Error(),
		})
		return
	}

	if err := DB.Delete(&siswa).Error; err != nil {
		c.JSON(http.StatusInternalServerError, APIResponse{
			Status:  "error",
			Message: "Gagal menghapus data siswa",
			Error:   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Siswa dengan ID %s berhasil dihapus", id),
	})
}

// scalarHandler menyajikan tampilan interaktif Scalar API Reference
func scalarHandler(c *gin.Context) {
	htmlContent := `<!doctype html>
<html>
  <head>
    <title>Siswa API Reference</title>
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <style>
      body {
        margin: 0;
      }
    </style>
  </head>
  <body>
    <script
      id="api-reference"
      data-url="/openapi.json"
      data-configuration='{"theme":"purple","layout":"modern","darkMode":true}'></script>
    <script src="https://cdn.jsdelivr.net/npm/@scalar/api-reference"></script>
  </body>
</html>`
	c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
}

// openAPIHandler menyajikan spesifikasi OpenAPI 3.0 dalam format JSON
func openAPIHandler(c *gin.Context) {
	c.Data(http.StatusOK, "application/json; charset=utf-8", openapiSpec)
}

func main() {
	initDB()

	r := gin.Default()

	// Dokumentasi API (Scalar & OpenAPI)
	r.GET("/openapi.json", openAPIHandler)
	r.GET("/scalar", scalarHandler)
	r.GET("/docs", scalarHandler)

	// CRUD Siswa Endpoints
	r.POST("/siswa", createSiswa)
	r.GET("/siswa", getAllSiswa)
	r.GET("/siswa/:id", getSiswaByID)
	r.PUT("/siswa/:id", updateSiswa)
	r.DELETE("/siswa/:id", deleteSiswa)

	fmt.Println("Aplikasi berjalan di http://localhost:8081")
	fmt.Println("Dokumentasi Scalar dapat diakses di http://localhost:8081/scalar atau http://localhost:8081/docs")

	if err := r.Run(":8081"); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

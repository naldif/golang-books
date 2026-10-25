package books

import (
	"errors"
	"net/http"
	"strconv"

	apperror "golang-buku/pkg/error"
	"golang-buku/pkg/response"

	"github.com/gin-gonic/gin"
)

type Controller struct {
	service Service
}

func NewController(service Service) *Controller {
	return &Controller{service: service}
}

func (ctrl *Controller) GetAll(c *gin.Context) {
	books, err := ctrl.service.GetAll()
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar buku berhasil diambil", books)
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	book, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, apperror.ErrBookNotFound) {
			response.Error(c, http.StatusNotFound, "Buku tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail buku berhasil diambil", book)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input CreateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)

	book, err := ctrl.service.Create(input, usernameStr)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidYear) {
			response.FieldErrors(c, http.StatusBadRequest, "Validasi data gagal", map[string][]string{
				"release_year": {"Tahun rilis harus di antara tahun 1980 dan 2024"},
			})
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Buku berhasil ditambahkan", book)
}

func (ctrl *Controller) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	var input UpdateBookInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)

	book, err := ctrl.service.Update(uint(id), input, usernameStr)
	if err != nil {
		if errors.Is(err, apperror.ErrBookNotFound) {
			response.Error(c, http.StatusNotFound, "Buku dengan ID tersebut tidak ditemukan", nil)
			return
		}
		if errors.Is(err, apperror.ErrInvalidYear) {
			response.FieldErrors(c, http.StatusBadRequest, "Validasi data gagal", map[string][]string{
				"release_year": {"Tahun rilis harus di antara tahun 1980 dan 2024"},
			})
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Buku berhasil diperbarui", book)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	err = ctrl.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, apperror.ErrBookNotFound) {
			response.Error(c, http.StatusNotFound, "Buku dengan ID tersebut tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Buku berhasil dihapus", nil)
}

package categories

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
	categories, err := ctrl.service.GetAll()
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar kategori berhasil diambil", categories)
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	category, err := ctrl.service.GetByID(uint(id))
	if err != nil {
		if errors.Is(err, apperror.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "Kategori tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail kategori berhasil diambil", category)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input CreateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)

	category, err := ctrl.service.Create(input, usernameStr)
	if err != nil {
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "Kategori berhasil ditambahkan", category)
}

func (ctrl *Controller) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	var input UpdateCategoryInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)

	category, err := ctrl.service.Update(uint(id), input, usernameStr)
	if err != nil {
		if errors.Is(err, apperror.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "Kategori dengan ID tersebut tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Kategori berhasil diperbarui", category)
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
		if errors.Is(err, apperror.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "Kategori dengan ID tersebut tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Kategori berhasil dihapus", nil)
}

func (ctrl *Controller) GetBooksByCategoryID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID tidak valid", nil)
		return
	}

	books, err := ctrl.service.GetBooksByCategoryID(uint(id))
	if err != nil {
		if errors.Is(err, apperror.ErrCategoryNotFound) {
			response.Error(c, http.StatusNotFound, "Kategori tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Daftar buku berdasarkan kategori berhasil diambil", books)
}

package users

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
	users, err := ctrl.service.FindAll()
	if err != nil {
		response.InternalServerError(c, err)
		return
	}
	response.Success(c, http.StatusOK, "Daftar user berhasil diambil", users)
}

func (ctrl *Controller) GetByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	user, err := ctrl.service.FindByID(uint(id))
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Detail user berhasil diambil", user)
}

func (ctrl *Controller) Create(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)
	if usernameStr == "" {
		usernameStr = "system"
	}

	user, err := ctrl.service.Create(input, usernameStr)
	if err != nil {
		if errors.Is(err, apperror.ErrUserExists) {
			response.FieldErrors(c, http.StatusBadRequest, "Validasi data gagal", map[string][]string{
				"username": {"Username sudah digunakan"},
			})
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusCreated, "User berhasil dibuat", user)
}

func (ctrl *Controller) Update(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	username, _ := c.Get("username")
	usernameStr, _ := username.(string)
	if usernameStr == "" {
		usernameStr = "system"
	}

	user, err := ctrl.service.Update(uint(id), input, usernameStr)
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User tidak ditemukan", nil)
			return
		}
		if errors.Is(err, apperror.ErrUserExists) {
			response.FieldErrors(c, http.StatusBadRequest, "Validasi data gagal", map[string][]string{
				"username": {"Username sudah digunakan oleh akun lain"},
			})
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User berhasil diperbarui", user)
}

func (ctrl *Controller) Delete(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "ID user tidak valid", nil)
		return
	}

	err = ctrl.service.Delete(uint(id))
	if err != nil {
		if errors.Is(err, apperror.ErrUserNotFound) {
			response.Error(c, http.StatusNotFound, "User tidak ditemukan", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "User berhasil dihapus", nil)
}

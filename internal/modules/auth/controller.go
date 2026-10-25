package auth

import (
	"errors"
	"net/http"

	"golang-buku/internal/modules/users"
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

func (ctrl *Controller) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	user, err := ctrl.service.Register(input)
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

	res := users.UserResponse{
		ID:         user.ID,
		Username:   user.Username,
		CreatedAt:  user.CreatedAt,
		CreatedBy:  user.CreatedBy,
		ModifiedAt: user.ModifiedAt,
		ModifiedBy: user.ModifiedBy,
	}

	response.Success(c, http.StatusCreated, "Registrasi berhasil", res)
}

func (ctrl *Controller) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.ValidationError(c, err)
		return
	}

	token, err := ctrl.service.Login(input)
	if err != nil {
		if errors.Is(err, apperror.ErrInvalidLogin) {
			response.Error(c, http.StatusUnauthorized, "Username atau password salah", nil)
			return
		}
		response.InternalServerError(c, err)
		return
	}

	response.Success(c, http.StatusOK, "Login berhasil", LoginResponse{Token: token})
}

func (ctrl *Controller) Logout(c *gin.Context) {
	response.Success(c, http.StatusOK, "Logout berhasil. Silakan hapus token autentikasi pada sisi klien.", nil)
}

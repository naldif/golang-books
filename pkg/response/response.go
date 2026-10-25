package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func init() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterTagNameFunc(func(fld reflect.StructField) string {
			name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
			if name == "-" {
				return ""
			}
			return name
		})
	}
}

type SuccessResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

type ErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Errors  any    `json:"errors"`
}

func Success(c *gin.Context, statusCode int, message string, data any) {
	c.PureJSON(statusCode, SuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string, err any) {
	var errVal any = err
	if errStr, ok := err.(string); ok && errStr == "" {
		errVal = nil
	}
	c.PureJSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  errVal,
	})
}

func ValidationError(c *gin.Context, err error) {
	fieldErrors := FormatValidationErrors(err)
	c.PureJSON(http.StatusBadRequest, ErrorResponse{
		Success: false,
		Message: "Validasi data gagal",
		Errors:  fieldErrors,
	})
}

func FieldErrors(c *gin.Context, statusCode int, message string, fieldErrors map[string][]string) {
	if message == "" {
		message = "Validasi data gagal"
	}
	c.PureJSON(statusCode, ErrorResponse{
		Success: false,
		Message: message,
		Errors:  fieldErrors,
	})
}

func InternalServerError(c *gin.Context, err error) {
	c.PureJSON(http.StatusInternalServerError, ErrorResponse{
		Success: false,
		Message: "Terjadi kesalahan pada server",
		Errors:  nil,
	})
}

func FormatValidationErrors(err error) map[string][]string {
	res := make(map[string][]string)

	var ve validator.ValidationErrors
	if errors.As(err, &ve) {
		for _, fe := range ve {
			field := fe.Field()
			if field == "" {
				field = strings.ToLower(fe.StructField())
			}

			var msg string
			switch fe.Tag() {
			case "required":
				msg = fmt.Sprintf("Field '%s' wajib diisi", field)
			case "min":
				if fe.Type().Kind() == reflect.String {
					msg = fmt.Sprintf("Field '%s' minimal %s karakter", field, fe.Param())
				} else {
					msg = fmt.Sprintf("Field '%s' minimal %s", field, fe.Param())
				}
			case "max":
				if fe.Type().Kind() == reflect.String {
					msg = fmt.Sprintf("Field '%s' maksimal %s karakter", field, fe.Param())
				} else {
					msg = fmt.Sprintf("Field '%s' maksimal %s", field, fe.Param())
				}
			case "email":
				msg = fmt.Sprintf("Field '%s' harus berupa email yang valid", field)
			default:
				msg = fmt.Sprintf("Field '%s' tidak valid", field)
			}

			res[field] = append(res[field], msg)
		}
		return res
	}

	var ute *json.UnmarshalTypeError
	if errors.As(err, &ute) {
		field := ute.Field
		if field == "" {
			field = "payload"
		}
		res[field] = append(res[field], fmt.Sprintf("Tipe data field '%s' tidak sesuai, diharapkan %s", field, ute.Type.String()))
		return res
	}

	res["payload"] = []string{"Format data JSON tidak valid"}
	return res
}

package users

import (
	"golang-buku/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller) {
	usersGroup := r.Group("/users")
	usersGroup.Use(middleware.AuthMiddleware())
	{
		usersGroup.GET("", ctrl.GetAll)
		usersGroup.POST("", ctrl.Create)
		usersGroup.GET("/:id", ctrl.GetByID)
		usersGroup.PUT("/:id", ctrl.Update)
		usersGroup.PATCH("/:id", ctrl.Update)
		usersGroup.DELETE("/:id", ctrl.Delete)
	}
}

package auth

import (
	"golang-buku/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller) {
	authGroup := r.Group("/auth")
	{
		authGroup.POST("/register", ctrl.Register)
		authGroup.POST("/login", ctrl.Login)
		authGroup.POST("/logout", middleware.AuthMiddleware(), ctrl.Logout)
	}
}

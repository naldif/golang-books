package books

import (
	"golang-buku/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller) {
	booksGroup := r.Group("/books")
	booksGroup.Use(middleware.AuthMiddleware())
	{
		booksGroup.GET("", ctrl.GetAll)
		booksGroup.POST("", ctrl.Create)
		booksGroup.GET("/:id", ctrl.GetByID)
		booksGroup.PUT("/:id", ctrl.Update)
		booksGroup.PATCH("/:id", ctrl.Update)
		booksGroup.DELETE("/:id", ctrl.Delete)
	}
}

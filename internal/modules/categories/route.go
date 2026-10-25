package categories

import (
	"golang-buku/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.RouterGroup, ctrl *Controller) {
	categoriesGroup := r.Group("/categories")
	categoriesGroup.Use(middleware.AuthMiddleware())
	{
		categoriesGroup.GET("", ctrl.GetAll)
		categoriesGroup.POST("", ctrl.Create)
		categoriesGroup.GET("/:id", ctrl.GetByID)
		categoriesGroup.PUT("/:id", ctrl.Update)
		categoriesGroup.PATCH("/:id", ctrl.Update)
		categoriesGroup.DELETE("/:id", ctrl.Delete)
		categoriesGroup.GET("/:id/books", ctrl.GetBooksByCategoryID)
	}
}

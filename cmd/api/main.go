package main

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"os"

	"golang-buku/internal/config"
	"golang-buku/internal/middleware"
	"golang-buku/internal/modules/auth"
	"golang-buku/internal/modules/books"
	"golang-buku/internal/modules/categories"
	"golang-buku/internal/modules/users"

	"github.com/gin-gonic/gin"
)

//go:embed openapi.json
var openapiSpec []byte

func main() {
	config.LoadEnv()

	db := config.InitDB()

	authRepo := auth.NewRepository(db)
	userRepo := users.NewRepository(db)
	categoryRepo := categories.NewRepository(db)
	bookRepo := books.NewRepository(db)

	authService := auth.NewService(authRepo)
	userService := users.NewService(userRepo)
	categoryService := categories.NewService(categoryRepo)
	bookService := books.NewService(bookRepo)

	authCtrl := auth.NewController(authService)
	userCtrl := users.NewController(userService)
	categoryCtrl := categories.NewController(categoryService)
	bookCtrl := books.NewController(bookService)

	r := gin.Default()

	r.Use(middleware.CORSMiddleware())

	swaggerHandler := func(c *gin.Context) {
		htmlContent := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Book Backend API - Swagger UI</title>
  <link rel="stylesheet" type="text/css" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-32x32.png" sizes="32x32" />
  <link rel="icon" type="image/png" href="https://unpkg.com/swagger-ui-dist@5/favicon-16x16.png" sizes="16x16" />
  <style>
    html {
      box-sizing: border-box;
      overflow: -moz-scrollbars-vertical;
      overflow-y: scroll;
    }
    *, *:before, *:after {
      box-sizing: inherit;
    }
    body {
      margin: 0;
      background: #fafafa;
    }
  </style>
</head>
<body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js" charset="UTF-8"></script>
  <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-standalone-preset.js" charset="UTF-8"></script>
  <script>
    window.onload = function() {
      const ui = SwaggerUIBundle({
        url: "/openapi.json",
        dom_id: '#swagger-ui',
        deepLinking: true,
        presets: [
          SwaggerUIBundle.presets.apis,
          SwaggerUIStandalonePreset
        ],
        plugins: [
          SwaggerUIBundle.plugins.DownloadUrl
        ],
        layout: "StandaloneLayout"
      });
      window.ui = ui;
    };
  </script>
</body>
</html>`
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(htmlContent))
	}

	r.GET("/", swaggerHandler)
	r.GET("/docs", swaggerHandler)
	r.GET("/swagger", swaggerHandler)

	r.GET("/openapi.json", func(c *gin.Context) {
		c.Data(http.StatusOK, "application/json; charset=utf-8", openapiSpec)
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"success": true,
			"message": "Book Backend API is running",
			"data": gin.H{
				"name":    "Book Backend API",
				"status":  "healthy",
				"version": "1.0.0",
			},
		})
	})

	apiGroup := r.Group("/api")
	{
		auth.RegisterRoutes(apiGroup, authCtrl)
		users.RegisterRoutes(apiGroup, userCtrl)
		categories.RegisterRoutes(apiGroup, categoryCtrl)
		books.RegisterRoutes(apiGroup, bookCtrl)

		apiGroup.POST("/users/register", authCtrl.Register)
		apiGroup.POST("/users/login", authCtrl.Login)
		apiGroup.POST("/users/logout", middleware.AuthMiddleware(), authCtrl.Logout)
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("Server berjalan di http://localhost:%s\n", port)
	fmt.Printf("Dokumentasi Swagger dapat diakses di http://localhost:%s atau http://localhost:%s/docs\n", port, port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
	}
}

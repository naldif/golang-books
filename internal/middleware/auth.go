package middleware

import (
	"net/http"
	"strings"

	"golang-buku/pkg/response"
	"golang-buku/utils"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Token autentikasi tidak ditemukan (wajib: Bearer <token>)", nil)
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" || strings.TrimSpace(parts[1]) == "" {
			response.Error(c, http.StatusUnauthorized, "Token autentikasi tidak ditemukan (wajib: Bearer <token>)", nil)
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := utils.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, http.StatusUnauthorized, "Token tidak valid atau telah kedaluwarsa", nil)
			c.Abort()
			return
		}

		c.Set("userID", claims.UserID)
		c.Set("username", claims.Username)

		c.Next()
	}
}

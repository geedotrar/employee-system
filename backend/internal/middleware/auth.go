package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"main.go/internal/auth"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			ctx.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")
		if token == authHeader {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			ctx.Abort()
			return
		}

		_, err := auth.ValidateJWT(token)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"main.go/internal/auth"
	"main.go/internal/models"
)

func AuthMiddleware(db *gorm.DB) gin.HandlerFunc {
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

		var blacklisted models.BlackListedToken
		if db == nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Database connection is nil"})
			ctx.Abort()
			return
		}

		if err := db.Where("token = ?", token).First(&blacklisted).Error; err == nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Token is blacklisted"})
			ctx.Abort()
			return
		} else if err != gorm.ErrRecordNotFound {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error checking blacklisted token", "details": err.Error()})
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

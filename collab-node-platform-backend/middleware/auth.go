package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"collab-node-platform-backend/utils"
	"collab-node-platform-backend/config"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" || !strings.HasPrefix(header, "Bearer ") {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "未携带有效Token"})
			return
		}
		token := strings.TrimPrefix(header, "Bearer ")
		claims, err := utils.ValidateJWT(token, config.AppConfig.JwtSecret)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token无效或已过期"})
			return
		}
		c.Set("userId", claims.UserID)
		c.Next()
	}
}

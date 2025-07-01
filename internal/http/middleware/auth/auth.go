package auth

import (
	"net/http"
	"qooked/internal/auth/jwt"
	"strings"

	"github.com/gin-gonic/gin"
)

func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// try to extract the authorization header
		authHeader := c.GetHeader("Authorization")
		// if header missing return 401 unauthorized
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is required"})
			return
		}

		// Bearer token format: "Bearer <token>"
		tokenString := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		if tokenString == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Bearer token is required"})
			return
		}

		// Validate the JWT token
		claims, err := jwt.ValidateJWT(tokenString)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			return
		}

		tokenUsername, ok := claims["username"].(string)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token payload"})
			return
		}

		routeUser := c.Param("username")
		if routeUser != tokenUsername {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "User mismatch: unauthorized access"})
			return
		}

		// Store the claims in the context
		c.Set("username", tokenUsername)
		c.Set("claims", claims)

		// Continue to the next middleware/handler
		c.Next()
	}
}

package middleware

import (
	"strings"

	"movie-ticket/internal/models"
	appErrors "movie-ticket/internal/utils/errors"
	"movie-ticket/internal/utils/jwt"
	"movie-ticket/internal/utils/response"

	"github.com/gin-gonic/gin"
)

const (
	CtxUserIDKey = "user_id"
	CtxUserEmail = "user_email"
	CtxUserRole  = "user_role"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, appErrors.Unauthorized("Authorization header is required", nil))
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && strings.ToLower(parts[0]) == "bearer") {
			response.Error(c, appErrors.Unauthorized("Authorization header format must be Bearer {token}", nil))
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwt.ValidateToken(tokenString)
		if err != nil {
			response.Error(c, appErrors.Unauthorized("Invalid or expired token: "+err.Error(), err))
			c.Abort()
			return
		}

		// Inject user claims into Gin context
		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxUserEmail, claims.Email)
		c.Set(CtxUserRole, claims.Role)

		c.Next()
	}
}

func RequireRole(roles ...models.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		val, exists := c.Get(CtxUserRole)
		if !exists {
			response.Error(c, appErrors.Unauthorized("Authentication required", nil))
			c.Abort()
			return
		}

		userRole, ok := val.(models.Role)
		if !ok {
			response.Error(c, appErrors.Forbidden("Invalid user role", nil))
			c.Abort()
			return
		}

		allowed := false
		for _, r := range roles {
			if userRole == r {
				allowed = true
				break
			}
		}

		if !allowed {
			response.Error(c, appErrors.Forbidden("Access denied: insufficient permissions", nil))
			c.Abort()
			return
		}

		c.Next()
	}
}

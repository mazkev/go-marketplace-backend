package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"go-market/internal/domain"
	"go-market/pkg/jwt"
	"go-market/pkg/response"
)

const (
	CtxUserIDKey = "user_id"
	CtxEmailKey  = "email"
	CtxRoleKey   = "role"
)

func AuthMiddleware(jwtService jwt.JWTService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Unauthorized(c, "Authorization header is required")
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			response.Unauthorized(c, "Invalid Authorization header format. Format must be 'Bearer <token>'")
			c.Abort()
			return
		}

		tokenString := parts[1]
		claims, err := jwtService.ValidateToken(tokenString)
		if err != nil {
			response.Unauthorized(c, "Invalid or expired token")
			c.Abort()
			return
		}

		c.Set(CtxUserIDKey, claims.UserID)
		c.Set(CtxEmailKey, claims.Email)
		c.Set(CtxRoleKey, claims.Role)
		c.Next()
	}
}

func RequireRole(allowedRoles ...domain.Role) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleVal, exists := c.Get(CtxRoleKey)
		if !exists {
			response.Forbidden(c, "Access forbidden: no role specified")
			c.Abort()
			return
		}

		userRole := domain.Role(roleVal.(string))
		for _, role := range allowedRoles {
			if userRole == role {
				c.Next()
				return
			}
		}

		response.Forbidden(c, "Access forbidden: insufficient permissions")
		c.Abort()
	}
}

func GetCurrentUserID(c *gin.Context) uint {
	id, exists := c.Get(CtxUserIDKey)
	if !exists {
		return 0
	}
	return id.(uint)
}

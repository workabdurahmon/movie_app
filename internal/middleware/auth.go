package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"errors"
)

type AuthMiddleware struct {
	jwtKey string
}

func NewAuthMiddleware(jwtKey string) *AuthMiddleware {
	return &AuthMiddleware{
		jwtKey: jwtKey,
	}
}

var (
	ErrMissingToken     = errors.New("missing authorization token")
	ErrInvalidAuthHeader = errors.New("invalid authorization header")
	ErrInvalidTokenFormat = errors.New("invalid token format")
	ErrUnauthorized      = errors.New("unauthorized")
)

func (m *AuthMiddleware) Authenticate() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": ErrMissingToken.Error()})
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": ErrInvalidAuthHeader.Error()})
			ctx.Abort()
			return
		}

		tokenString := parts[1]
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			return []byte(m.jwtKey), nil
		})

		if err != nil || !token.Valid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": ErrUnauthorized.Error()})
			ctx.Abort()
			return
		}

		claims, isValid := token.Claims.(jwt.MapClaims)
		if !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			ctx.Abort()
			return
		}

		userID, isValid := claims["user_id"].(float64)
		if !isValid {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID in token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", uint(userID))
		ctx.Next()
	}
}

package middleware

import (
	"billohub/config"
	"billohub/pkg/helper"
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWTMiddleware creates a gin middleware for JWT authentication.
func JWTMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get token from Authorization header
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.NewErrorResponse(helper.ErrToken.WithMessage("Authorization header is missing"), nil))
			return
		}

		// The token is expected to be in the format "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.NewErrorResponse(helper.ErrToken.WithMessage("Authorization header format must be Bearer {token}"), nil))
			return
		}
		tokenString := parts[1]

		// 2. Parse and validate the token
		cfg := config.GetConfig()
		jwtKey := []byte(cfg.JwtKey)

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("alg not validate")
			}
			return jwtKey, nil
		})

		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.NewErrorResponse(helper.ErrToken.WithMessage("Invalid token: "+err.Error()), nil))
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			// 3. Store user info in context for downstream handlers
			c.Set("username", claims["username"])
			c.Set("role", claims["role"])
			c.Next()
		} else {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.NewErrorResponse(helper.ErrToken.WithMessage("Invalid token claims"), nil))
		}
	}
}

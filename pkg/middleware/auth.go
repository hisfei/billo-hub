package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"billohub/pkg/helper"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
)

// ContextKey is a custom Context Key type to avoid key collisions.
type ContextKey string

const (
	UserIDKey      ContextKey = "userId"
	UserRoleKey    ContextKey = "userRole"
	IsCertifiedKey ContextKey = "isCertified"
	LevelKey       ContextKey = "level"
	StatusKey      ContextKey = "status"
	CreatedAtKey   ContextKey = "createdAt"
)

// AuthStruct defines the user authentication information structure stored in Redis.
type AuthStruct struct {
	Token       string `json:"token"`
	UserRole    string `json:"userRole"`
	IsCertified bool   `json:"isCertified"` // Recommended to use bool type
	Level       int    `json:"level"`       // Recommended to use int type
	Status      bool   `json:"status"`
	CreatedAt   string `json:"createdAt"`
}

// AdvancedAuthMiddleware returns a Gin middleware for advanced authentication.
// It relies on a user session whitelist stored in Redis.
// Prerequisite: X-UID and X-TOKEN must be verified by a signature in a preceding middleware.
func AdvancedAuthMiddleware(rdb *redis.Client) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Get authentication information from the header
		uidStr := c.GetHeader("X-UId")
		tokenStr := c.GetHeader("X-TOKEN")

		// Validate that X-UID is a valid number
		uid, err := strconv.ParseInt(uidStr, 10, 64)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, helper.NewErrorResponse(helper.ErrWrongParam.WithMessage("Invalid X-UID format"), nil))
			return
		}

		// 2. Get user authentication information from Redis
		var authStruct AuthStruct
		redisKey := fmt.Sprintf("auth:%d", uid)
		err = rdb.HGetAll(context.Background(), redisKey).Scan(&authStruct)

		// 3. Validate login status and token
		if err == redis.Nil || authStruct.Token != tokenStr {
			c.AbortWithStatusJSON(http.StatusUnauthorized, helper.NewErrorResponse(helper.ErrTokenExpire.WithMessage("Login status invalid or expired"), nil))
			return
		}
		if err != nil {
			// Other Redis errors
			c.AbortWithStatusJSON(http.StatusInternalServerError, helper.NewErrorResponse(helper.ErrInner.WithMessage("Redis error during auth"), nil))
			return
		}

		// 4. Check if the user is banned
		if !authStruct.Status {
			c.AbortWithStatusJSON(http.StatusForbidden, helper.NewErrorResponse(helper.ErrUserAuth.WithMessage("Account is banned"), nil))
			return
		}

		// 5. Store the "absolutely trusted" user information in the context
		c.Set(string(UserIDKey), uid)
		c.Set(string(UserRoleKey), authStruct.UserRole)
		c.Set(string(IsCertifiedKey), authStruct.IsCertified)
		c.Set(string(LevelKey), authStruct.Level)
		c.Set(string(StatusKey), authStruct.Status)
		c.Set(string(CreatedAtKey), authStruct.CreatedAt)

		c.Next()
	}
}

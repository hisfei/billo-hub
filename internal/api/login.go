package api

import (
	"billohub/config"
	"billohub/internal/model"
	"billohub/pkg/helper"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Login handles user authentication.
func (h *APIHandler) Login(c *gin.Context) {
	var loginData struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var res helper.APIResponse
	res.CodeDetail = helper.OK

	if err := c.ShouldBindJSON(&loginData); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		if h.DebugMode {
			res.Msg = err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// Authenticate user via AgentHub
	user, err := h.Hub.LoginUser(loginData.Username, loginData.Password)
	if err != nil {
		res.CodeDetail = helper.ErrUserAuth
		res.Msg = "Invalid username or password"
		if h.DebugMode {
			res.Msg += ": " + err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	// Generate JWT Token
	cfg := config.GetConfig()
	jwtKey := []byte(cfg.JwtKey)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"role":     "admin", // Assuming all users are 'admin' for now, can be extended
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		if h.DebugMode {
			res.Msg = "Failed to sign token: " + err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res.Body = map[string]interface{}{
		"token": tokenString,
		"role":  "admin",
	}
	c.JSON(http.StatusOK, res)
}

// ResetPassword handles user password reset requests.
func (h *APIHandler) ResetPassword(c *gin.Context) {
	var req model.ResetPasswordRequest
	var res helper.APIResponse
	res.CodeDetail = helper.OK

	if err := c.ShouldBindJSON(&req); err != nil {
		res.CodeDetail = helper.ErrWrongParam
		res.Msg = err.Error()
		c.JSON(http.StatusOK, res)
		return
	}

	// It's good practice to let the authenticated user only reset their own password.
	// We can get the username from the JWT token claims.
	usernameFromToken, exists := c.Get("username")
	if !exists || usernameFromToken.(string) != req.Username {
		res.CodeDetail = helper.ErrUserAuth
		res.Msg = "You can only reset your own password."
		c.JSON(http.StatusForbidden, res)
		return
	}

	err := h.Hub.ResetUserPassword(req.Username, req.OldPassword, req.NewPassword)
	if err != nil {
		res.CodeDetail = helper.ErrInner
		res.Msg = "Failed to reset password"
		if h.DebugMode {
			res.Msg += ": " + err.Error()
		}
		c.JSON(http.StatusOK, res)
		return
	}

	res.Msg = "Password reset successfully"
	c.JSON(http.StatusOK, res)
}

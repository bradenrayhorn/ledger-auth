package routing

import (
	"github.com/bradenrayhorn/ledger-auth/internal"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/gin-gonic/gin"
	"net/http"
)

type RegisterRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

type LoginRequest struct {
	Username string `form:"username" binding:"required"`
	Password string `form:"password" binding:"required"`
}

func Register(c *gin.Context) {
	var request RegisterRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(internal.MakeValidationError(err))
		return
	}

	if err := services.RegisterUser(request.Username, request.Password); err != nil {
		_ = c.Error(err)
	}
}

func Login(c *gin.Context) {
	var request LoginRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(internal.MakeValidationError(err))
		return
	}

	token, err := services.Login(request.Username, request.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})
}

func Me(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"id":       c.GetString("user_id"),
		"username": c.GetString("user_username"),
	})
}
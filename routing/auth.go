package routing

import (
	"github.com/bradenrayhorn/ledger-auth/internal"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/gin-gonic/gin"
)

type RegisterRequest struct {
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

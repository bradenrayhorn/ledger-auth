package routing

import (
	"errors"
	"net/mail"

	"github.com/bradenrayhorn/ledger-auth/internal"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/gin-gonic/gin"
)

type UpdateEmailRequest struct {
	Email string `form:"email"`
}

func UpdateEmail(c *gin.Context) {
	var request UpdateEmailRequest
	if err := c.ShouldBind(&request); err != nil {
		_ = c.Error(internal.MakeValidationError(err))
		return
	}

	if len(request.Email) > 0 {
		_, err := mail.ParseAddress(request.Email)
		if err != nil {
			_ = c.Error(internal.MakeValidationError(errors.New("invalid email")))
			return
		}
	}

	err := services.UpdateEmail(c.GetString("user_id"), request.Email)
	if err != nil {
		_ = c.Error(err)
		return
	}
}

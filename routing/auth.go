package routing

import (
	"context"
	"net/http"

	"github.com/bradenrayhorn/ledger-auth/database"
	"github.com/bradenrayhorn/ledger-auth/internal"
	"github.com/bradenrayhorn/ledger-auth/repositories"
	"github.com/bradenrayhorn/ledger-auth/services"
	"github.com/gin-gonic/gin"
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

	// add device id if needed
	cookie, err := c.Request.Cookie("device_id")
	var deviceID string
	if err != nil {
		newDeviceID, err := createDeviceCookie(c.Writer)
		if err != nil {
			_ = c.Error(err)
			return
		}
		deviceID = newDeviceID
	} else {
		deviceID = cookie.Value
	}

	// login user
	user, err := services.Login(request.Username, request.Password)
	if err != nil {
		_ = c.Error(err)
		return
	}

	deviceService := services.NewDeviceService(database.RDB)
	// check if device is recognized
	isRecognized, err := deviceService.DoesRecognizeDevice(context.Background(), user.ID, deviceID)
	if err != nil {
		_ = c.Error(err)
		return
	} else if !isRecognized {
		if err = deviceService.NotifyOfNewDevice(*user, c.ClientIP()); err != nil {
			_ = c.Error(err)
			return
		}
	}

	// recognize device
	err = deviceService.RecognizeDevice(context.Background(), user.ID, deviceID)
	if err != nil {
		_ = c.Error(err)
		return
	}

	err = createSession(c.Writer, user.ID, c.ClientIP(), c.Request.UserAgent())
	if err != nil {
		_ = c.Error(err)
		return
	}
}

func Logout(c *gin.Context) {
	err := deleteSession(c.Writer, c.GetString("session_id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
}

func Me(c *gin.Context) {
	user, err := repositories.GetUserByID(context.Background(), c.GetString("user_id"))
	if err != nil {
		_ = c.Error(err)
		return
	}

	var email *string
	if user.Email.Valid {
		email = &user.Email.String
	}

	c.IndentedJSON(http.StatusOK, map[string]interface{}{
		"id":    c.GetString("user_id"),
		"email": email,
	})
}

func RevokeSessions(c *gin.Context) {
	err := services.NewSessionService(database.RDB).DeleteActiveSessionsForUser(context.Background(), c.GetString("user_id"))
	if err != nil {
		_ = c.Error(err)
		return
	}
}

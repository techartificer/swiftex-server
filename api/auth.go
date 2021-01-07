package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/validators"
)

// RegisterAuthRoutes initialize all auth related routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/login/", login)
}

func login(ctx echo.Context) error {
	resp := response.Response{}
	_, err := validators.ValidateLogin(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid login request data"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	return resp.Send(ctx)
}

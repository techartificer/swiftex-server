package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/validators"
)

// RegisterAdminRoutes initialize all auth related routes
func RegisterAdminRoutes(endpoint *echo.Group) {
	endpoint.POST("/add", createAdmin, middlewares.JWTAuth())
}

func createAdmin(ctx echo.Context) error {
	resp := response.Response{}
	admin, err := validators.ValidateAddAdmin(ctx)
	if err != nil {
		logger.Errorln(err)
		resp.Title = "Invalid add admin request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	logger.Infoln(admin)
	return resp.Send(ctx)
}

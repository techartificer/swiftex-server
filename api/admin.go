package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/validators"
)

// RegisterAdminRoutes initialize all auth related routes
func RegisterAdminRoutes(endpoint *echo.Group) {
	// endpoint.POST("/add/", createAdmin, middlewares.JWTAuth())
	endpoint.POST("/add/", createAdmin)
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
	db := database.GetDB()
	logger.Infoln(admin)
	adminRepo := data.NewAdminRepo()
	if err := adminRepo.Create(db, admin); err != nil {
		logger.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Status = http.StatusCreated
	resp.Data = admin
	return resp.Send(ctx)
}

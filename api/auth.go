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
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAuthRoutes initialize all auth related routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/admin/login/", adminLogin)
}

func adminLogin(ctx echo.Context) error {
	resp := response.Response{}
	body, err := validators.ValidateLogin(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid login request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	admin, err := adminRepo.FindByUsername(db, body.Phone)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Admin not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.AdminNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = admin
	return resp.Send(ctx)
}

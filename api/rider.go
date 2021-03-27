package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/password"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/models"
	"github.com/techartificer/swiftex/validators"
)

func RegisterRiderRoutes(endpoint *echo.Group) {
	endpoint.POST("/create/", createRider, middlewares.JWTAuth(true))
	endpoint.GET("/", riders, middlewares.JWTAuth(true))
	endpoint.GET("/:hub/", ridersByHub, middlewares.JWTAuth(true))
}

func createRider(ctx echo.Context) error {
	resp := response.Response{}
	rider, err := validators.ValidateRiderReq(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid rider create request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	hash, err := password.HashPassword(rider.Password)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Password hash failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.PasswordHashFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	rider.Password = hash
	riderRepo := data.NewRiderRepo()
	db := database.GetDB()

	if err := riderRepo.Create(db, rider); err != nil {
		logger.Log.Errorln(err)
		if errors.IsMongoDupError(err) {
			resp.Title = "Rider already exist"
			resp.Status = http.StatusConflict
			resp.Code = codes.MerchantAlreadyExist
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = rider
	resp.Status = http.StatusCreated
	return resp.Send(ctx)
}

func ridersByHub(ctx echo.Context) error {
	resp := response.Response{}
	hub := ctx.Param("hub")
	db := database.GetDB()
	riderRepo := data.NewRiderRepo()
	riders, err := riderRepo.RidersByHub(db, hub)

	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = func() []models.Rider {
		if *riders == nil {
			return []models.Rider{}
		}
		return *riders
	}()
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func riders(ctx echo.Context) error {
	resp := response.Response{}
	db := database.GetDB()
	riderRepo := data.NewRiderRepo()
	lastID := ctx.QueryParam("lastId")

	riders, err := riderRepo.Riders(db, lastID)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = func() []models.Rider {
		if *riders == nil {
			return []models.Rider{}
		}
		return *riders
	}()
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

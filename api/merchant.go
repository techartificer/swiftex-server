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
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterMerchantRoutes(endpoint *echo.Group) {
	endpoint.POST("/register/", register)
	endpoint.GET("/is-available/:phone/", isUsernameAvilable)
}

func isUsernameAvilable(ctx echo.Context) error {
	resp := response.Response{}
	phone := ctx.Param("phone")
	merchantRepo := data.NewMerchantRepo()
	db := database.GetDB()
	_, err := merchantRepo.FindByPhone(db, phone)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp.Data = map[string]bool{"available": true}
			resp.Status = http.StatusOK
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = map[string]bool{"available": false}
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func register(ctx echo.Context) error {
	resp := response.Response{}
	merchant, err := validators.ValidateMerchantRegister(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid register request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	merchantRepo := data.NewMerchantRepo()

	hash, err := password.HashPassword(merchant.Password)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Password hash failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.PasswordHashFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	merchant.Password = hash
	if err := merchantRepo.Create(db, merchant); err != nil {
		logger.Log.Errorln(err)
		if errors.IsMongoDupError(err) {
			resp.Title = "Merchant already exist"
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
	resp.Data = merchant
	resp.Status = http.StatusCreated
	return resp.Send(ctx)
}

package api

import (
	"net/http"

	"github.com/gosimple/slug"
	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterShopRoutes(endpoint *echo.Group) {
	endpoint.POST("/create/", create, middlewares.JWTAuth())
}

func create(ctx echo.Context) error {
	resp := response.Response{}
	shop, err := validators.ValidateShopCreate(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid shop create request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	shopRepo := data.NewShopRepo()
	shop.Owner = ctx.Get(constants.UserID).(primitive.ObjectID)
	shop.ShopID = slug.Make(shop.Name)

	if err := shopRepo.Create(db, shop); err != nil {
		logger.Log.Errorln(err)
		if errors.IsMongoDupError(err) {
			resp.Title = "Shop already exist"
			resp.Status = http.StatusConflict
			resp.Code = codes.ShopAlreadyExist
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = shop
	resp.Status = http.StatusCreated
	return resp.Send(ctx)
}

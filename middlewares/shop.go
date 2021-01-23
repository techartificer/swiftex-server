package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func isModerator(shop *models.Shop, userID primitive.ObjectID) bool {
	if shop.Moderators == nil {
		return false
	}
	shop.Moderators = []primitive.ObjectID{primitive.NewObjectID()}
	for _, v := range shop.Moderators {
		if v == userID {
			return true
		}
	}
	return false
}

func HasShopAccess() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := response.Response{}
			r := ctx.Get(constants.Role).(string)
			role := constants.AdminRole(r)
			shopID := ctx.Param("shopId")
			db := database.GetDB()
			shopRepo := data.NewShopRepo()
			shop, err := shopRepo.ShopByID(db, shopID)
			if err != nil {
				logger.Log.Errorln(err)
				if err == mongo.ErrNoDocuments {
					resp.Title = "Shop not found"
					resp.Status = http.StatusNotFound
					resp.Code = codes.ShopNotFound
					resp.Errors = errors.NewError(err.Error())
					return resp.Send(ctx)
				}
				resp.Title = "Something went wrong"
				resp.Status = http.StatusInternalServerError
				resp.Code = codes.DatabaseQueryFailed
				resp.Errors = err
				return resp.Send(ctx)
			}
			isAdmin := role == constants.Admin || role == constants.SuperAdmin || role == constants.Moderator || role == constants.ZoneManager
			if isAdmin {
				ctx.Set("shop", shop)
				return next(ctx)
			}
			userID := ctx.Get(constants.UserID).(primitive.ObjectID)
			isModerator := isModerator(shop, userID)
			if shop.Owner != userID && !isModerator {
				resp.Title = "You don not have access"
				resp.Status = http.StatusForbidden
				resp.Code = codes.AccessDenied
				resp.Errors = err
				return resp.Send(ctx)
			}
			ctx.Set("shop", shop)
			return next(ctx)
		}
	}
}

func IsShopOwner() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := response.Response{}
			r := ctx.Get(constants.Role).(string)
			role := constants.AdminRole(r)

			shopID := ctx.Param("shopId")
			db := database.GetDB()
			shopRepo := data.NewShopRepo()
			shop, err := shopRepo.ShopByID(db, shopID)
			if err != nil {
				logger.Log.Errorln(err)
				if err == mongo.ErrNoDocuments {
					resp.Title = "Shop not found"
					resp.Status = http.StatusNotFound
					resp.Code = codes.ShopNotFound
					resp.Errors = errors.NewError(err.Error())
					return resp.Send(ctx)
				}
				resp.Title = "Something went wrong"
				resp.Status = http.StatusInternalServerError
				resp.Code = codes.DatabaseQueryFailed
				resp.Errors = err
				return resp.Send(ctx)
			}
			isAdmin := role == constants.Admin || role == constants.SuperAdmin
			if isAdmin {
				ctx.Set("shop", shop)
				return next(ctx)
			}
			userID := ctx.Get(constants.UserID).(primitive.ObjectID)
			if shop.Owner != userID {
				resp.Title = "You don not have access"
				resp.Status = http.StatusForbidden
				resp.Code = codes.AccessDenied
				resp.Errors = err
				return resp.Send(ctx)
			}
			ctx.Set("shop", shop)
			return next(ctx)
		}
	}
}

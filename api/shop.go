package api

import (
	"net/http"
	"strconv"
	"time"

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
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterShopRoutes(endpoint *echo.Group) {
	endpoint.POST("/create/", shopCreate, middlewares.JWTAuth(false))
	endpoint.GET("/myshops/", myShops, middlewares.JWTAuth(false))
	endpoint.GET("/all-shops/", allShops, middlewares.JWTAuth(true))
	endpoint.GET("/id/:shopId/", shopByID, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.PATCH("/id/:shopId/", updateShop, middlewares.JWTAuth(false), middlewares.IsShopOwner())
	endpoint.GET("/search/", searchShop, middlewares.JWTAuth(true))
	endpoint.GET("/dashboard/:shopId/", dashboard, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/all-shops-name/", allShopsName, middlewares.JWTAuth(true))
}

func dashboard(ctx echo.Context) error {
	resp := response.Response{}
	shopID := ctx.Param("shopId")
	startDate, endDate := ctx.QueryParam("startDate"), ctx.QueryParam("endData")
	var tms, tme time.Time
	if startDate != "" && endDate != "" {
		std, err := strconv.ParseInt(startDate, 10, 64) // startDate
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid timestamp"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.SomethingWentWrong
			resp.Errors = err
			return resp.Send(ctx)
		}
		tms = time.Unix(std/1000, 0) //std => startDate

		end, err := strconv.ParseInt(endDate, 10, 64)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid timestamp"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.SomethingWentWrong
			resp.Errors = err
			return resp.Send(ctx)
		}
		tme = time.Unix(end/1000, 0)
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()

	dashboard, err := orderRepo.Dashboard(db, shopID, &tms, &tme)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = dashboard
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func allShopsName(ctx echo.Context) error {
	resp := response.Response{}
	shopRepo := data.NewShopRepo()
	db := database.GetDB()
	shops, err := shopRepo.AllShopsName(db)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = shops
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func searchShop(ctx echo.Context) error {
	resp := response.Response{}
	name, phone := ctx.QueryParam("name"), ctx.QueryParam("phone")
	query := make(bson.M)
	if name != "" {
		query["name"] = primitive.Regex{Pattern: name, Options: "i"}
	}
	if phone != "" {
		query["phone"] = primitive.Regex{Pattern: phone, Options: ""}
	}
	db := database.GetDB()
	shopRepo := data.NewShopRepo()
	shops, err := shopRepo.Search(db, query)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = shops
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func shopCreate(ctx echo.Context) error {
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
	trx, err := shopRepo.Create(db, shop)
	if err != nil {
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
	resp.Data = map[string]interface{}{
		"shop":        shop,
		"transaction": trx,
	}
	resp.Status = http.StatusCreated
	return resp.Send(ctx)
}

func allShops(ctx echo.Context) error {
	resp := response.Response{}
	lastID, limit := ctx.QueryParam("lastId"), ctx.QueryParam("limit")
	var limitNum int64
	if limit != "" {
		ln, err := strconv.Atoi(limit)
		limitNum = int64(ln)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Errors = err
			resp.Title = "Invalid limit"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InvalidLimit
			return resp.Send(ctx)
		}
	}
	db := database.GetDB()
	shopRepo := data.NewShopRepo()
	shops, err := shopRepo.Shops(db, lastID, limitNum)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
	}
	resp.Data = shops
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func updateShop(ctx echo.Context) error {
	resp := response.Response{}
	shop, err := validators.ValidateShopUpdate(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid shop update request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	if len(shop.Name) > 0 {
		shop.ShopID = slug.Make(shop.Name)
	}
	shopID := ctx.Param("shopId")
	db := database.GetDB()
	shopRepo := data.NewShopRepo()
	updatedShop, err := shopRepo.UpdateShopByID(db, shopID, shop)
	if err != nil {
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
	resp.Status = http.StatusOK
	resp.Data = updatedShop
	return resp.Send(ctx)
}

func myShops(ctx echo.Context) error {
	resp := response.Response{}
	db := database.GetDB()
	shopRepo := data.NewShopRepo()
	ownerID := ctx.Get(constants.UserID).(primitive.ObjectID)
	shops, err := shopRepo.ShopsByOwnerId(db, ownerID)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Can not fetch data"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = shops
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func shopByID(ctx echo.Context) error {
	resp := response.Response{}
	shop := ctx.Get("shop")
	resp.Data = shop
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

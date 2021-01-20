package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/random"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func RegisterOrderRoutes(endpoint *echo.Group) {
	endpoint.POST("/create/", orderCreate, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/all/:shopId/", orders, middlewares.JWTAuth(false), middlewares.HasShopAccess())
}

func orders(ctx echo.Context) error {
	resp := response.Response{}
	shopID := ctx.Param("shopId")
	lastID, isDelivered := ctx.QueryParam("lastId"), ctx.QueryParam("isDelivered")
	_shopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid shop ID"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = codes.InvalidMongoID
		resp.Errors = err
		return resp.Send(ctx)
	}
	query := make(bson.M)
	query["shopId"] = _shopID
	if lastID != "" {
		id, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid order ID"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InvalidMongoID
			resp.Errors = err
			return resp.Send(ctx)
		}
		query["_id"] = bson.M{"$gt": id}
	}
	date := time.Date(2020, 12, 31, 0, 0, 0, 0, time.UTC)
	if isDelivered != "" {
		query["deliverdAt"] = bson.M{"$gt": date}
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	orders, err := orderRepo.Orders(db, query)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.SomethingWentWrong
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Status = http.StatusOK
	resp.Data = orders
	return resp.Send(ctx)
}

func orderCreate(ctx echo.Context) error {
	resp := response.Response{}
	order, err := validators.ValidateOrderCreate(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid order create request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidShopCreateData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	tid, err := random.GenerateRandomString(constants.TrackIDSize)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.SomethingWentWrong
		resp.Errors = err
		return resp.Send(ctx)
	}
	order.TrackID = tid
	if err := orderRepo.Create(db, order); err != nil {
		logger.Log.Errorln(err)
		if errors.IsMongoDupError(err) {
			resp.Title = "Order track Id already exist"
			resp.Status = http.StatusConflict
			resp.Code = codes.OrderAlreadyExist
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = order
	resp.Status = http.StatusCreated
	return resp.Send(ctx)
}

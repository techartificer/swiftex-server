package api

import (
	"log"
	"net/http"
	"strconv"
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
	"github.com/techartificer/swiftex/models"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterOrderRoutes(endpoint *echo.Group) {
	endpoint.GET("/", ordersAdmin, middlewares.JWTAuth(true))
	endpoint.POST("/create/:shopId/", orderCreate, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/all/:shopId/", orders, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.PATCH("/id/:orderId/shopId/:shopId/", updateOrder, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.PATCH("/add/order-status/:orderId/", addOrderStatus, middlewares.JWTAuth(true)) // TODO: Delivery boy access
	endpoint.PATCH("/cancel/id/:orderId/shopId/:shopId/", cancelOrder, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/id/:orderId/shopId/:shopId/", orderByID, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/track/:trackId/", trackOrder)
	endpoint.POST("/assign-rider/", assignRider, middlewares.JWTAuth(true))
	endpoint.GET("/riders-parcel/:riderId/", ridersParcel, middlewares.RiderJWTAuth())
	endpoint.POST("/deliver/:orderId/", deliverParcel, middlewares.RiderJWTAuth())
}

func deliverParcel(ctx echo.Context) error {
	resp := response.Response{}
	trxHistory, err := validators.ValidateOrderDeliver(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid delivery request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidTransactionData
		resp.Errors = err
		return resp.Send(ctx)
	}
	trxRepo := data.NewTransactionRepo()
	db := database.GetDB()

	result, err := trxRepo.AddTrxHistory(db, trxHistory)
	if err != nil {
		logger.Log.Errorln(err)
		if mongo.ErrNoDocuments == err {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.OrderAlreadyDelevired) {
			resp.Title = "Order already delivered"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.OrderAlreadyDelevired
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.TransactionNotFound) {
			resp.Title = "Transaction not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.TransactionNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.OrderNotAcceptedYet) {
			resp.Title = "Order not accepted yet"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.OrderNotAcceptedYet
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = result
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func ridersParcel(ctx echo.Context) error {
	resp := response.Response{}
	riderID := ctx.Param("riderId")
	lastID := ctx.QueryParam("lastId")

	db := database.GetDB()
	ridersParcelRepo := data.NewRiderParcelRepo()
	orders, err := ridersParcelRepo.ParcelsByRiderId(db, riderID, lastID)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = func() []bson.M {
		if *orders == nil {
			return []bson.M{}
		}
		return *orders
	}()
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func assignRider(ctx echo.Context) error {
	resp := response.Response{}
	body, err := validators.ValidateRiderParcelCreate(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid assign parcel request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidAssignParcelData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	riderParcelRepo := data.NewRiderParcelRepo()
	order, err := riderParcelRepo.Create(db, body)
	if err != nil {
		logger.Log.Errorln(err)
		if mongo.ErrNoDocuments == err {
			resp.Title = "Order not exist"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.OrderAlreadyInTransit) {
			resp.Title = "Order already assigned to a rider"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.OrderAlreadyInTransit
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.OrderAlreadyDelevired) {
			resp.Title = "Order already delivered"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.OrderAlreadyDelevired
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Status = http.StatusCreated
	resp.Data = map[string]interface{}{
		"order":       order,
		"riderParcel": body,
	}
	return resp.Send(ctx)
}

func ordersAdmin(ctx echo.Context) error {
	resp := response.Response{}
	lastID, startDate, endDate, shopID := ctx.QueryParam("lastId"), ctx.QueryParam("startDate"), ctx.QueryParam("endDate"), ctx.QueryParam("shopId")
	trackID, phone := ctx.QueryParam("trackId"), ctx.QueryParam("phone")
	query := make(bson.M)
	if shopID != "" {
		_shopID, err := primitive.ObjectIDFromHex(shopID)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid shop ID"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InvalidMongoID
			resp.Errors = err
			return resp.Send(ctx)
		}
		query["shopId"] = _shopID
	}
	if lastID != "" {
		id, err := primitive.ObjectIDFromHex(lastID)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid last order ID"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InvalidMongoID
			resp.Errors = err
			return resp.Send(ctx)
		}
		query["_id"] = bson.M{"$lt": id}
	}
	if phone != "" {
		query["recipientPhone"] = primitive.Regex{Pattern: phone, Options: ""}
	}
	if trackID != "" {
		query["trackId"] = primitive.Regex{Pattern: trackID, Options: ""}
	}
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
		tms := time.Unix(std/1000, 0) //std => startDate
		end, err := strconv.ParseInt(endDate, 10, 64)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid timestamp"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.SomethingWentWrong
			resp.Errors = err
			return resp.Send(ctx)
		}
		tme := time.Unix(end/1000, 0)
		query["$and"] = []bson.M{{"createdAt": bson.M{"$gte": tms}}, {"createdAt": bson.M{"$lte": tme}}}
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

func trackOrder(ctx echo.Context) error {
	resp := response.Response{}
	trackID := ctx.Param("trackId")
	orderRepo := data.NewOrderRepo()
	db := database.GetDB()
	order, err := orderRepo.TrackOrder(db, trackID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = order.Status
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func updateOrder(ctx echo.Context) error {
	resp := response.Response{}
	orderID, shopID := ctx.Param("orderId"), ctx.Param("shopId")
	body, err := validators.UpdateOrder(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid order update request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidOrderUpdateData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()

	order, err := orderRepo.OrderByID(db, orderID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	if order.DeliveredAt != nil || order.IsPicked || order.IsCancelled {
		resp.Title = "You can not update parcel"
		resp.Status = http.StatusLocked
		resp.Code = codes.OrderNotUpdateAble
		resp.Errors = errors.NewError("Parcel status is not allowing to update")
		return resp.Send(ctx)
	}
	updatedOrder, err := orderRepo.UpdateOrder(db, body, orderID, shopID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = updatedOrder
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func cancelOrder(ctx echo.Context) error {
	resp := response.Response{}
	orderID, shopID := ctx.Param("orderId"), ctx.Param("shopId")

	cancelled := constants.Cancelled
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	order := &models.Order{
		CurrentStatus: &cancelled,
		IsCancelled:   true,
		UpdatedAt:     time.Now().UTC(),
	}
	updatedOrder, err := orderRepo.UpdateOrder(db, order, orderID, shopID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Order not found"
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
	resp.Data = updatedOrder
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func orderByID(ctx echo.Context) error {
	resp := response.Response{}
	orderID, shopID := ctx.Param("orderId"), ctx.Param("shopId")

	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	order, err := orderRepo.OrderByID(db, orderID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}

	if order.ShopID.Hex() != shopID {
		resp.Title = "You don't have access"
		resp.Status = http.StatusForbidden
		resp.Code = codes.AccessDenied
		resp.Errors = errors.NewError(err.Error())
		return resp.Send(ctx)
	}

	resp.Data = order
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func addOrderStatus(ctx echo.Context) error {
	/*
		TODO: check order status
		TODO: If accepted can not update
	*/

	resp := response.Response{}
	orderID := ctx.Param("orderId")
	body, err := validators.UpdateOrderStatus(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid order update request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidOrderStatusUpdateData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	order, err := orderRepo.OrderByID(db, orderID)
	if err != nil {
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}

	if order.DeliveredAt != nil {
		resp.Title = "Order already delivered"
		resp.Status = http.StatusUnprocessableEntity
		resp.Code = codes.OrderAlreadyDelevired
		return resp.Send(ctx)
	}

	orderStatus, err := orderRepo.AddOrderStatus(db, body, orderID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Order not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.OrderNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = orderStatus
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func orders(ctx echo.Context) error {
	resp := response.Response{}
	shopID := ctx.Param("shopId")
	lastID, startDate, endDate := ctx.QueryParam("lastId"), ctx.QueryParam("startDate"), ctx.QueryParam("endDate")
	trackID, phone := ctx.QueryParam("trackId"), ctx.QueryParam("phone")

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
			resp.Title = "Invalid last order ID"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InvalidMongoID
			resp.Errors = err
			return resp.Send(ctx)
		}
		query["_id"] = bson.M{"$lt": id}
	}
	if phone != "" {
		query["recipientPhone"] = primitive.Regex{Pattern: phone, Options: ""}
	}
	if trackID != "" {
		query["trackId"] = primitive.Regex{Pattern: trackID, Options: ""}
	}
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
		tms := time.Unix(std/1000, 0) //std => startDate
		end, err := strconv.ParseInt(endDate, 10, 64)
		if err != nil {
			logger.Log.Errorln(err)
			resp.Title = "Invalid timestamp"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.SomethingWentWrong
			resp.Errors = err
			return resp.Send(ctx)
		}
		tme := time.Unix(end/1000, 0)
		query["$and"] = []bson.M{{"createdAt": bson.M{"$gte": tms}}, {"createdAt": bson.M{"$lte": tme}}}
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	orders, err := orderRepo.Orders(db, query)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
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
	shopID := ctx.Param("shopId")
	_shopID, err := primitive.ObjectIDFromHex(shopID)
	if err != nil {
		log.Println(err)
		resp.Title = "Invalid shop ID"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidShopCreateData
		resp.Errors = err
		return resp.Send(ctx)
	}
	order.ShopID = _shopID
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

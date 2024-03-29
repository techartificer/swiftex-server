package api

import (
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/charge"
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

type orderError struct {
	Error   string             `json:"error"`
	OrderID primitive.ObjectID `json:"orderId"`
}

func RegisterOrderRoutes(endpoint *echo.Group) {
	endpoint.GET("/", ordersAdmin, middlewares.JWTAuth(true))
	endpoint.POST("/create/:shopId/", orderCreate, middlewares.JWTAuth(false), middlewares.HasShopAccess(), middlewares.ShopByID())
	endpoint.GET("/all/:shopId/", orders, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.PATCH("/id/:orderId/shopId/:shopId/", updateOrder, middlewares.JWTAuth(false), middlewares.HasShopAccess(), middlewares.ShopByID())
	endpoint.PATCH("/add/order-status/:orderId/", addOrderStatus, middlewares.JWTAuth(true)) // TODO: Delivery boy access
	endpoint.PATCH("/cancel/id/:orderId/shopId/:shopId/", cancelOrder, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/id/:orderId/shopId/:shopId/", orderByID, middlewares.JWTAuth(false), middlewares.HasShopAccess())
	endpoint.GET("/track/:trackId/", trackOrder)
	endpoint.POST("/assign-rider/", assignRider, middlewares.JWTAuth(true))
	endpoint.GET("/riders-parcel/:riderId/", ridersParcel, middlewares.RiderJWTAuth())
	endpoint.POST("/deliver/:orderId/", deliverParcel, middlewares.RiderJWTAuth())
	endpoint.PATCH("/change/status/", changeStatus, middlewares.JWTAuth(true))
	endpoint.POST("/create/:shopId/multiples/", createMultipleOrder, middlewares.JWTAuth(false), middlewares.HasShopAccess(), middlewares.ShopByID())
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
		if err.Error() == string(codes.ShopNotFound) {
			resp.Title = "Shop not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.ShopNotFound
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
	trackID, phone, deliveryZone := ctx.QueryParam("trackId"), ctx.QueryParam("phone"), ctx.QueryParam("deliveryZone")
	query := make(bson.M)
	if deliveryZone != "" {
		query["recipientArea"] = primitive.Regex{Pattern: deliveryZone, Options: "i"}
	}
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
	orderID := ctx.Param("orderId")
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
	shop := ctx.Get("shop").(models.Shop)
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
	body.Charge = charge.Calculate(body.Weight, body.DeliveryType, body.RecipientCity, shop.DeliveryCharge)
	if order.DeliveredAt != nil || order.IsPicked || order.IsCancelled {
		resp.Title = "You can not update parcel"
		resp.Status = http.StatusLocked
		resp.Code = codes.OrderNotUpdateAble
		resp.Errors = errors.NewError("Parcel status is not allowing to update")
		return resp.Send(ctx)
	}
	updatedOrder, err := orderRepo.UpdateOrder(db, body, orderID, shop.ID.Hex())
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
	trackID, phone, deliveryZone := ctx.QueryParam("trackId"), ctx.QueryParam("phone"), ctx.QueryParam("deliveryZone")

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
	if deliveryZone != "" {
		query["recipientArea"] = primitive.Regex{Pattern: deliveryZone, Options: ""}
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
	db := database.GetDB()
	shop := ctx.Get("shop").(models.Shop)
	if err != nil {
		logger.Log.Errorln(err)
		if mongo.ErrNoDocuments == err {
			resp.Title = "Shop not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.ShopNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	order.Charge = charge.Calculate(order.Weight, order.DeliveryType, order.RecipientCity, shop.DeliveryCharge)
	order.ShopID = shop.ID
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

func changeStatus(ctx echo.Context) error {
	resp := response.Response{}
	body, err := validators.OrderChangeStatus(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid order update request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidOrderStatusUpdateData
		resp.Errors = err
		return resp.Send(ctx)
	}

	if body.Status == constants.Delivered {
		resp.Title = "Invalid order status change request"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidOrderStatusUpdateData
		resp.Errors = err
		return resp.Send(ctx)
	}

	var wg sync.WaitGroup
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	errChan := make(chan orderError, len(body.OrderIDs))
	ordersChan := make(chan models.Order, len(body.OrderIDs))

	for _, orderId := range body.OrderIDs {
		wg.Add(1)
		go func(w *sync.WaitGroup, oid primitive.ObjectID) {
			defer w.Done()
			order, err := orderRepo.OrderByID(db, oid.Hex())
			if err != nil {
				errChan <- orderError{err.Error(), oid}
				return
			}
			if order.DeliveredAt != nil {
				errChan <- orderError{"Order already delivered", oid}
				return
			}
			orderStatus := models.OrderStatus{
				ID:            primitive.NewObjectID(),
				Text:          body.Text,
				DeleveryBoyID: body.DeleveryBoyID,
				AdminID:       body.AdminID,
				Status:        body.Status,
				Time:          time.Now().UTC(),
			}
			order, err = orderRepo.AddOrderStatus(db, &orderStatus, oid.Hex())
			if err != nil {
				errChan <- orderError{err.Error(), oid}
			} else {
				ordersChan <- *order
			}
		}(&wg, orderId)
	}
	wg.Wait()
	close(errChan)
	close(ordersChan)

	hasError := false
	var errs []orderError
	for err := range errChan {
		hasError = true
		logger.Log.Errorln(err)
		errs = append(errs, err)
	}
	if hasError {
		resp.Errors = errors.NewError("Update status not successfull")
		resp.Data = errs
		resp.Status = http.StatusInternalServerError
		resp.Title = "Status update unsuccessful"
		resp.Code = codes.DatabaseQueryFailed
		return resp.Send(ctx)
	}
	var orders []models.Order
	for order := range ordersChan {
		orders = append(orders, order)
	}
	resp.Data = orders
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func createMultipleOrder(ctx echo.Context) error {
	resp := response.Response{}
	orders, err := validators.ValidateMultipleOrderCreate(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid order create request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidShopCreateData
		resp.Errors = err
		return resp.Send(ctx)
	}
	shop := ctx.Get("shop").(models.Shop)
	var os []interface{}
	for i := 0; i < len(orders); i++ {
		order := &orders[i]
		order.ShopID = shop.ID
		order.Charge = charge.Calculate(order.Weight, order.DeliveryType, order.RecipientCity, shop.DeliveryCharge)
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
		os = append(os, order)
	}
	db := database.GetDB()
	orderRepo := data.NewOrderRepo()
	if err := orderRepo.CreateMultiple(db, os); err != nil {
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
	resp.Data = orders
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterTransactionRoutes(endpoint *echo.Group) {
	endpoint.GET("/shopId/:shopId/", transactionByShopId, middlewares.JWTAuth(false), middlewares.IsShopOwner())
	endpoint.PATCH("/generate-trx-code/:shopId/", generateTrxCode, middlewares.JWTAuth(false), middlewares.IsShopOwnerStrict())
	endpoint.GET("/cash-out-requests/", cashOutRequests, middlewares.JWTAuth(true))
	endpoint.PATCH("/cash-out/:trxId/", makeCashOut, middlewares.JWTAuth(true))
}

func makeCashOut(ctx echo.Context) error {
	resp := response.Response{}
	trxID := ctx.Param("trxId")
	body, err := validators.ValidateCahsOutReq(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid cash out request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidCashOutData
		resp.Errors = err
		return resp.Send(ctx)
	}
	trxRepo := data.NewTransactionRepo()
	db := database.GetDB()

	userID := ctx.Get(constants.UserID).(primitive.ObjectID)

	trx, err := trxRepo.CashOut(db, userID, trxID, body.TrxCode)
	if err != nil {
		logger.Log.Errorln(err)
		if mongo.ErrNoDocuments == err {
			resp.Title = "Transaction not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.TransactionNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.InsufficientBalance) {
			resp.Title = "Insufficient balance"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InsufficientBalance
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.InvalidTrxCode) {
			resp.Title = "Invalid Trx code"
			resp.Status = http.StatusForbidden
			resp.Code = codes.InvalidTrxCode
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.TrxCodeExpired) {
			resp.Title = "Trx code expired"
			resp.Status = http.StatusForbidden
			resp.Code = codes.TrxCodeExpired
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = trx
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func cashOutRequests(ctx echo.Context) error {
	resp := response.Response{}
	lastID := ctx.QueryParam("lastId")
	db := database.GetDB()
	trxRepo := data.NewTransactionRepo()
	result, err := trxRepo.CashOutRequests(db, lastID)
	if err != nil {
		logger.Log.Errorln(err)
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

func transactionByShopId(ctx echo.Context) error {
	resp := response.Response{}
	shopID := ctx.Param("shopId")
	db := database.GetDB()
	trxRepo := data.NewTransactionRepo()

	result, err := trxRepo.TransactionByShopId(db, shopID)
	if err != nil {
		logger.Log.Errorln(err)
		if mongo.ErrNoDocuments == err {
			resp.Title = "Transaction not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.TransactionNotFound
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

func generateTrxCode(ctx echo.Context) error {
	resp := response.Response{}
	shopID := ctx.Param("shopId")
	body, err := validators.ValidateGenerateTrxCodeReq(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid Trx Code generate request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidGenTrxCode
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	trxRepo := data.NewTransactionRepo()

	trxCode, err := trxRepo.GenerateTrxCode(db, body.Amount, shopID)
	if err != nil {
		logger.Log.Errorln(err)
		if mongo.ErrNoDocuments == err {
			resp.Title = "Transaction not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.TransactionNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		if err.Error() == string(codes.InsufficientBalance) {
			resp.Title = "Insufficient balance"
			resp.Status = http.StatusUnprocessableEntity
			resp.Code = codes.InsufficientBalance
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = map[string]string{"trxCode": *trxCode}
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

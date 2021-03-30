package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterTransactionRoutes(endpoint *echo.Group) {
	endpoint.GET("/shopId/:shopId/", transactionByShopId)
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

package middlewares

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/lib/jwt"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func JWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := response.Response{}
			claims, _, err := jwt.ExtractAndValidateToken(ctx)
			if err != nil {
				logger.Log.Errorln(err)
				resp.Status = http.StatusUnauthorized
				resp.Code = codes.InvalidAuthorizationToken
				resp.Title = "Unauthorized request"
				resp.Errors = err
				return resp.Send(ctx)
			}
			userID, err := primitive.ObjectIDFromHex(claims.UserID)
			if err != nil {
				logger.Log.Errorln(err)
				resp.Title = "Something went wrong"
				resp.Status = http.StatusInternalServerError
				resp.Code = codes.SomethingWentWrong
				resp.Errors = err
				return resp.Send(ctx)
			}
			ctx.Set(constants.UserID, userID)
			ctx.Set(constants.Role, claims.Audience)
			ctx.Set(constants.Phone, claims.Phone)
			return next(ctx)
		}
	}
}

func IsSuperAdmin() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := response.Response{}
			role := ctx.Get(constants.Role)
			if role != string(constants.SuperAdmin) {
				resp.Status = http.StatusForbidden
				resp.Code = codes.NotSuperAdmin
				resp.Title = "You are not super admin"
				return resp.Send(ctx)
			}
			return next(ctx)
		}
	}
}
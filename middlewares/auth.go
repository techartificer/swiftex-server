package middlewares

import (
	"net/http"

	jwtGo "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/lib/jwt"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func setHeader(resp response.Response, ctx echo.Context, claims *jwt.Claims) error {
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
	return nil
}

func RiderJWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := response.Response{}
			claims, _, err := jwt.ExtractAndValidateToken(ctx)
			if err != nil {
				logger.Log.Errorln(err)
				ve, _ := err.(*jwtGo.ValidationError)
				resp.Code = codes.InvalidAuthorizationToken
				if ve.Errors == jwtGo.ValidationErrorExpired {
					resp.Code = codes.JWTExpired
				}
				resp.Status = http.StatusUnauthorized
				resp.Title = err.Error()
				resp.Errors = err
				return resp.Send(ctx)
			}
			if claims.AccountType != constants.AdminType && claims.AccountType != constants.RiderType {
				resp.Status = http.StatusForbidden
				resp.Code = codes.InvalidAccountType
				resp.Title = "You are not allowed"
				return resp.Send(ctx)
			}
			if err := setHeader(resp, ctx, claims); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}

func JWTAuth(isAdmin bool) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			resp := response.Response{}
			claims, _, err := jwt.ExtractAndValidateToken(ctx)
			if err != nil {
				logger.Log.Errorln(err)
				ve, _ := err.(*jwtGo.ValidationError)
				resp.Status = http.StatusUnauthorized
				resp.Code = codes.InvalidAuthorizationToken
				if ve.Errors == jwtGo.ValidationErrorExpired {
					resp.Code = codes.JWTExpired
				}
				resp.Title = err.Error()
				resp.Errors = err
				return resp.Send(ctx)
			}
			if isAdmin && claims.AccountType != constants.AdminType {
				resp.Status = http.StatusForbidden
				resp.Code = codes.InvalidAccountType
				resp.Title = "You are not allowed"
				return resp.Send(ctx)
			}
			if err := setHeader(resp, ctx, claims); err != nil {
				return err
			}
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

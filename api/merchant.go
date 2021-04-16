package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/firebase"
	"github.com/techartificer/swiftex/lib/jwt"
	"github.com/techartificer/swiftex/lib/password"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/models"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func RegisterMerchantRoutes(endpoint *echo.Group) {
	endpoint.POST("/register/", register)
	endpoint.GET("/is-available/:phone/", isUsernameAvilable)
	endpoint.GET("/", allMerchants, middlewares.JWTAuth(true))
	endpoint.PATCH("/forgot-password/", forgotPassword)
}

func isUsernameAvilable(ctx echo.Context) error {
	resp := response.Response{}
	phone := ctx.Param("phone")
	merchantRepo := data.NewMerchantRepo()
	db := database.GetDB()
	_, err := merchantRepo.FindByPhone(db, phone)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Data = map[string]bool{"available": true}
			resp.Status = http.StatusOK
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = map[string]bool{"available": false}
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func register(ctx echo.Context) error {
	resp := response.Response{}
	merchant, err := validators.ValidateMerchantRegister(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid register request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}

	token := ctx.Request().Header.Get("FirebaseToken")
	if err := firebase.ValidateToken(token, merchant.Phone); err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Phone number is not verified"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.PhoneNumberNotVerified
		resp.Errors = err
		return resp.Send(ctx)
	}

	db := database.GetDB()
	merchantRepo := data.NewMerchantRepo()

	hash, err := password.HashPassword(merchant.Password)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Password hash failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.PasswordHashFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	merchant.Password = hash
	if err := merchantRepo.Create(db, merchant); err != nil {
		logger.Log.Errorln(err)
		if errors.IsMongoDupError(err) {
			resp.Title = "Merchant already exist"
			resp.Status = http.StatusConflict
			resp.Code = codes.MerchantAlreadyExist
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	signedToken, err := jwt.BuildJWTToken(merchant.Phone, constants.ShopOwner, merchant.ID.Hex(), constants.MerchantType)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.UserLoginFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	sess := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       merchant.ID,
		RefreshToken: jwt.NewRefresToken(merchant.ID),
		AccessToken:  signedToken,
		CreatedAt:    time.Now().UTC(),
		ExpiresOn:    time.Now().Add(time.Minute * time.Duration(config.GetJWT().RefreshTTL)),
	}
	sessRepo := data.NewSessionRepo()
	if err = sessRepo.CreateSession(db, sess); err != nil {
		logger.Log.Errorln(err)
		resp.Title = "User login failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	result := map[string]interface{}{
		"accessToken":  sess.AccessToken,
		"refreshToken": sess.RefreshToken,
		"expiresOn":    sess.ExpiresOn,
		"permission":   "Owner",
	}
	resp.Data = result
	resp.Status = http.StatusCreated
	return resp.Send(ctx)
}

func allMerchants(ctx echo.Context) error {
	resp := response.Response{}
	lastID := ctx.QueryParam("lastId")
	merchantRepo := data.NewMerchantRepo()
	db := database.GetDB()

	merchants, err := merchantRepo.Merchants(db, lastID)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = merchants
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func forgotPassword(ctx echo.Context) error {
	resp := response.Response{}
	body, err := validators.ValidateForgotPassword(ctx)

	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidForgotPassData
		resp.Errors = err
		return resp.Send(ctx)
	}

	if err := firebase.ValidateToken(body.Token, body.Phone); err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Phone number is not verified"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.PhoneNumberNotVerified
		resp.Errors = err
		return resp.Send(ctx)
	}

	hash, err := password.HashPassword(body.Password)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Password hash failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.PasswordHashFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	merchantRepo := data.NewMerchantRepo()
	db := database.GetDB()
	merchant, err := merchantRepo.UpdateByPhone(db, body.Phone, &models.Merchant{Password: hash})
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "You are not registered"
			resp.Status = http.StatusNotFound
			resp.Code = codes.AdminNotFound
			resp.Errors = err
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	signedToken, err := jwt.BuildJWTToken(merchant.Phone, constants.ShopOwner, merchant.ID.Hex(), constants.MerchantType)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.UserLoginFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	sess := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       merchant.ID,
		RefreshToken: jwt.NewRefresToken(merchant.ID),
		AccessToken:  signedToken,
		CreatedAt:    time.Now().UTC(),
		ExpiresOn:    time.Now().Add(time.Minute * time.Duration(config.GetJWT().RefreshTTL)),
	}
	sessRepo := data.NewSessionRepo()
	if err = sessRepo.CreateSession(db, sess); err != nil {
		logger.Log.Errorln(err)
		resp.Title = "User login failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	result := map[string]interface{}{
		"accessToken":  sess.AccessToken,
		"refreshToken": sess.RefreshToken,
		"expiresOn":    sess.ExpiresOn,
		"permission":   "Owner",
	}
	resp.Data = result
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

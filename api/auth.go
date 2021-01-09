package api

import (
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/config"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/jwt"
	"github.com/techartificer/swiftex/lib/password"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/models"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAuthRoutes initialize all auth related routes
func RegisterAuthRoutes(endpoint *echo.Group) {
	endpoint.POST("/admin/login/", adminLogin)
	endpoint.DELETE("/logout/", logout)
	endpoint.PATCH("/refresh-token/", refreshToken)
}

func adminLogin(ctx echo.Context) error {
	resp := response.Response{}
	body, err := validators.ValidateLogin(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid login request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	admin, err := adminRepo.FindByUsername(db, body.Phone)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Admin not found"
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
	if ok := password.CheckPasswordHash(body.Password, admin.Password); !ok {
		resp.Title = "Phone number or password incorrect"
		resp.Status = http.StatusUnauthorized
		resp.Code = codes.InvalidLoginCredential
		resp.Errors = err
		return resp.Send(ctx)
	}
	signedToken, err := jwt.BuildJWTToken(admin.Phone, string(admin.Role), admin.ID.Hex())
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
		UserID:       admin.ID,
		RefreshToken: jwt.NewRefresToken(admin.ID),
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
		"expireOn":     sess.ExpiresOn,
		"permission":   admin.Role,
	}
	resp.Status = http.StatusOK
	resp.Data = result
	return resp.Send(ctx)
}

func logout(ctx echo.Context) error {
	resp := response.Response{}
	token, err := jwt.ParseRefreshToken(ctx)
	if err != nil {
		resp.Title = "You are already logged out"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.BearerTokenGiven
		resp.Errors = err
		return resp.Send(ctx)
	}
	sessionRepo := data.NewSessionRepo()
	db := database.GetDB()
	if err := sessionRepo.Logout(db, token); err != nil {
		if err == mongo.ErrNoDocuments {
			resp.Title = "You are already logged out"
			resp.Status = http.StatusNotFound
			resp.Code = codes.RefreshTokenNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Logout failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Status = http.StatusOK
	resp.Title = "Logout successful"
	return resp.Send(ctx)
}

func refreshToken(ctx echo.Context) error {
	resp := response.Response{}
	token, err := jwt.ParseRefreshToken(ctx)
	if err != nil {
		resp.Title = "Token parsing failed"
		resp.Errors = err
		resp.Status = http.StatusBadRequest
		resp.Code = codes.UserSignUpDataInvalid
		return resp.Send(ctx)
	}
	db := database.GetDB()
	sessionRepo := data.NewSessionRepo()
	splittedToken := strings.Split(token, ".")
	userID, err := primitive.ObjectIDFromHex(splittedToken[1])
	if err != nil {
		resp.Title = "Invalid refresh token"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.TokenRefreshFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	adminRepo := data.NewAdminRepo()
	admin, err := adminRepo.FindByID(db, userID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Admin not found"
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
	accessToken, err := jwt.BuildJWTToken(admin.Phone, string(admin.Role), admin.ID.Hex())
	if err != nil {
		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.TokenRefreshFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	sess, err := sessionRepo.UpdateSession(db, token, accessToken, userID)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			resp.Title = "You are logged out"
			resp.Status = http.StatusNotFound
			resp.Code = codes.RefreshTokenNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Token refresh failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = sess
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

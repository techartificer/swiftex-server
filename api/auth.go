package api

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
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
		logger.Log.Infoln(err)

		resp.Title = "Failed to sign auth token"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.UserLoginFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	sess := &models.Session{
		ID:           primitive.NewObjectID(),
		UserID:       admin.ID,
		RefreshToken: jwt.NewRefresToken(),
		AccessToken:  signedToken,
		CreatedAt:    time.Now().UTC(),
		ExpiresOn:    time.Now().Add(time.Hour * 1),
	}
	sessRepo := data.NewSessionRepo()
	if err = sessRepo.CreateSession(db, sess); err != nil {
		logger.Log.Infoln(err)
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

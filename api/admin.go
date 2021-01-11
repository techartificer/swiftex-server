package api

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/data"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/errors"
	"github.com/techartificer/swiftex/lib/password"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/techartificer/swiftex/logger"
	"github.com/techartificer/swiftex/middlewares"
	"github.com/techartificer/swiftex/validators"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// RegisterAdminRoutes initialize all auth related routes
func RegisterAdminRoutes(endpoint *echo.Group) {
	endpoint.POST("/add/", createAdmin, middlewares.JWTAuth(), middlewares.IsSuperAdmin())
	endpoint.PATCH("/update/:adminId/", updateAdmin, middlewares.JWTAuth(), middlewares.IsSuperAdmin())
	endpoint.GET("/all/", allAdmins, middlewares.JWTAuth(), middlewares.IsSuperAdmin())
	endpoint.GET("/profile/", profile, middlewares.JWTAuth())
}

func createAdmin(ctx echo.Context) error {
	resp := response.Response{}
	admin, err := validators.ValidateAddAdmin(ctx)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid add admin request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	hash, err := password.HashPassword(admin.Password)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Password hash failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.PasswordHashFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	admin.Password = hash
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	if err := adminRepo.Create(db, admin); err != nil {
		logger.Log.Errorln(err)
		if errors.IsMongoDupError(err) {
			resp.Title = "Admin already exist"
			resp.Status = http.StatusConflict
			resp.Code = codes.AdminAlreadyExist
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
	resp.Data = admin
	return resp.Send(ctx)
}

func updateAdmin(ctx echo.Context) error {
	resp := response.Response{}
	ID := ctx.Param("adminId")
	body, err := validators.ValidateAdminUpdate(ctx)
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Invalid admin update request data"
		resp.Status = http.StatusBadRequest
		resp.Code = codes.InvalidRegisterData
		resp.Errors = err
		return resp.Send(ctx)
	}
	admin, err := adminRepo.UpdateAdminByID(db, body, ID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Admin not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.AdminNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Admin update failed"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	sessionRepo := data.NewSessionRepo()
	if body.Status == constants.Deactive {
		if _, err := sessionRepo.RemoveSessionsByUserID(db, ID); err != nil {
			resp.Title = "Admin update failed"
			resp.Status = http.StatusInternalServerError
			resp.Code = codes.DatabaseQueryFailed
			resp.Errors = err
			return resp.Send(ctx)
		}
	}
	resp.Data = admin
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func allAdmins(ctx echo.Context) error {
	resp := response.Response{}
	db := database.GetDB()
	adminRepo := data.NewAdminRepo()
	admins, err := adminRepo.AdminList(db)
	if err != nil {
		logger.Log.Errorln(err)
		resp.Title = "Can not fetch data"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = admins
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

func profile(ctx echo.Context) error {
	resp := response.Response{}
	db := database.GetDB()
	userID := ctx.Get(constants.UserID).(primitive.ObjectID)
	adminRepo := data.NewAdminRepo()
	admin, err := adminRepo.FindByID(db, userID)
	if err != nil {
		logger.Log.Errorln(err)
		if err == mongo.ErrNoDocuments {
			resp.Title = "Admin not found"
			resp.Status = http.StatusNotFound
			resp.Code = codes.AdminNotFound
			resp.Errors = errors.NewError(err.Error())
			return resp.Send(ctx)
		}
		resp.Title = "Something went wrong"
		resp.Status = http.StatusInternalServerError
		resp.Code = codes.DatabaseQueryFailed
		resp.Errors = err
		return resp.Send(ctx)
	}
	resp.Data = admin
	resp.Status = http.StatusOK
	return resp.Send(ctx)
}

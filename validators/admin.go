package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReqAdminAdd holds admin add request data
type ReqAdminAdd struct {
	Phone    string              `json:"phone,omitempty" validate:"required"`
	Name     string              `json:"name,omitempty" validate:"required"`
	Email    string              `json:"email,omitempty" validate:"required,email"`
	Role     constants.AdminRole `json:"role,omitempty" validate:"required,isValidRole"`
	Password string              `json:"password,omitempty" validate:"required,min=6,max=26"`
}

func isValidRole(fl validator.FieldLevel) bool {
	isFound := false
	for _, v := range constants.Roles {
		if string(v) == fl.Field().String() {
			isFound = true
		}
	}
	return isFound
}

// ValidateAddAdmin returns admin or error
func ValidateAddAdmin(ctx echo.Context) (*models.Admin, error) {
	v.RegisterValidation("isValidRole", isValidRole)
	body := ReqAdminAdd{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	admin := &models.Admin{
		ID:        primitive.NewObjectID(),
		Name:      body.Name,
		Email:     body.Email,
		Phone:     body.Phone,
		Password:  body.Password,
		Role:      body.Role,
		Status:    constants.Active,
		CreatedAt: time.Now().UTC(),
	}
	return admin, nil
}

// ReqAdminUpdate holds status update data request data
type ReqAdminUpdate struct {
	Phone    string              `json:"phone,omitempty"  bson:"phone,omitempty"`
	Name     string              `json:"name,omitempty" bson:"name,omitempty"`
	Email    *string             `json:"email,omitempty" validate:"omitempty,email" bson:"email,omitempty"`
	Role     constants.AdminRole `json:"role,omitempty" validate:"omitempty,isValidRole" bson:"role,omitempty"`
	Status   string              `json:"status,omitempty" validate:"omitempty,isValidStatus" bson:"status,omitempty"`
	Password string              `json:"password,omitempty" validate:"omitempty" bson:"password,omitempty"`
}

func isValidStatus(fl validator.FieldLevel) bool {
	isFound := false
	for _, v := range constants.AllStatus {
		if string(v) == fl.Field().String() {
			isFound = true
		}
	}
	return isFound
}

func ValidateAdminUpdate(ctx echo.Context) (*ReqAdminUpdate, error) {
	v.RegisterValidation("isValidRole", isValidRole)
	v.RegisterValidation("isValidStatus", isValidStatus)
	body := ReqAdminUpdate{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	return &body, nil
}

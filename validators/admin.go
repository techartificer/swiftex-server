package validators

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// ReqLogin holds login request data
type ReqAdminAdd struct {
	Phone    string              `json:"phone,omitempty" validate:"required"`
	Name     string              `json:"name,omitempty" validate:"required"`
	Email    string              `json:"email,omitempty" validate:"required,email"`
	Role     constants.AdminRole `json:"role,omitempty" validate:"required"`
	Password string              `json:"password,omitempty" validate:"required,min=6,max=26"`
}

// ValidateLogin returns request body or error
func ValidateAddAdmin(ctx echo.Context) (*models.Admin, error) {
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

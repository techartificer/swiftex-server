package validators

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type MerchantRegisterReq struct {
	Phone    string `json:"phone,omitempty" validate:"required"`
	Name     string `json:"name,omitempty" validate:"required"`
	Email    string `json:"email,omitempty" validate:"required,email"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
}

func ValidateMerchantRegister(ctx echo.Context) (*models.Merchant, error) {
	body := MerchantRegisterReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	merchant := &models.Merchant{
		ID:        primitive.NewObjectID(),
		Name:      body.Name,
		Password:  body.Password,
		Email:     body.Email,
		Phone:     body.Phone,
		Status:    constants.Active,
		CreatedAt: time.Now().UTC(),
	}
	return merchant, nil
}

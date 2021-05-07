package validators

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ShopCreateReq struct {
	Phone         string `json:"phone,omitempty" validate:"required"`
	Name          string `json:"name,omitempty" validate:"required,min=3,max=30"`
	Email         string `json:"email,omitempty" validate:"required,email"`
	Address       string `json:"address,omitempty" validate:"omitempty"`
	PickupAddress string `json:"pickupAddress,omitempty" validate:"required"`
	DeliveryZone  string `json:"deliveryZone,omitempty" validate:"omitempty"`
	FBPage        string `json:"fbPage,omitempty" validate:"required"`
	PickupArea    string `json:"pickupArea,omitempty" validate:"required"`
}

func ValidateShopCreate(ctx echo.Context) (*models.Shop, error) {
	body := ShopCreateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	shop := &models.Shop{
		ID:             primitive.NewObjectID(),
		Name:           body.Name,
		Email:          body.Email,
		Phone:          body.Phone,
		Address:        body.Address,
		PickupAddress:  body.PickupAddress,
		PickupArea:     body.PickupArea,
		FBPage:         body.FBPage,
		DeliveryZone:   body.DeliveryZone,
		Status:         constants.Active,
		DeliveryCharge: constants.DeliveryCharge,
		COD:            constants.CodCharge,
		CreatedAt:      time.Now().UTC(),
	}
	return shop, nil
}

type ShopUpdateReq struct {
	Phone         string `json:"phone,omitempty"`
	Name          string `json:"name,omitempty" validate:"omitempty,min=3,max=30"`
	Email         string `json:"email,omitempty" validate:"omitempty,email"`
	Address       string `json:"address,omitempty"`
	PickupAddress string `json:"pickupAddress,omitempty"`
	DeliveryZone  string `json:"deliveryZone,omitempty"`
	FBPage        string `json:"fbPage,omitempty"`
	PickupArea    string `json:"pickupArea,omitempty"`
}

func ValidateShopUpdate(ctx echo.Context) (*models.Shop, error) {
	body := ShopUpdateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	shop := &models.Shop{
		Name:          body.Name,
		Email:         body.Email,
		Phone:         body.Phone,
		Address:       body.Address,
		PickupAddress: body.PickupAddress,
		PickupArea:    body.PickupArea,
		FBPage:        body.FBPage,
		DeliveryZone:  body.DeliveryZone,
		UpdatedAt:     time.Now().UTC(),
	}
	return shop, nil
}

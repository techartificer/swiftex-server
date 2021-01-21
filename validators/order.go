package validators

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderCreateReq struct {
	ShopID                primitive.ObjectID `validate:"required" json:"shopId"`
	DeliveryBoy           primitive.ObjectID `validate:"omitempty" json:"deliveryBoy"`
	ShopModeratorID       primitive.ObjectID `validate:"omitempty" json:"shopModerator"`
	MerchantID            primitive.ObjectID `validate:"omitempty" json:"merchant"`
	RecipientName         string             `validate:"required" json:"recipientName"`
	RecipientPhone        string             `validate:"required" json:"recipientPhone"`
	RecipientCity         string             `validate:"required" json:"recipientCity"`
	RecipientThana        string             `validate:"required" json:"recipientThana"`
	RecipientArea         string             `validate:"required" json:"recipientArea"`
	RecipientZip          string             `validate:"required" json:"recipientZip"`
	RecipientAddress      string             `validate:"required" json:"recipientAddress"`
	PackageCode           string             `validate:"omitempty" json:"packageCode"`
	PaymentStatus         string             `validate:"required" json:"paymentStatus"`
	Price                 float64            `validate:"required" json:"price"`
	ParcelType            string             `validate:"required" json:"parcelType"`
	RequestedDeliveryTime time.Time          `validate:"omitempty" json:"requestedDeliveryTime"`
	PickAddress           string             `validate:"required" json:"pickAddress"`
	PickHub               string             `validate:"required" json:"pickHub"`
	Comments              string             `validate:"omitempty,max=300" json:"comments"`
	NumberOfItems         int                `validate:"required" json:"numberOfItems"`
	DeliveryType          string             `validate:"required" json:"deliveryType"`
}

func ValidateOrderCreate(ctx echo.Context) (*models.Order, error) {
	body := OrderCreateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	order := &models.Order{
		ID:                    primitive.NewObjectID(),
		ShopID:                body.ShopID,
		DeliveryBoy:           body.DeliveryBoy,
		ShopModeratorID:       body.ShopModeratorID,
		MerchantID:            body.MerchantID,
		RecipientName:         body.RecipientName,
		RecipientPhone:        body.RecipientPhone,
		RecipientCity:         body.RecipientCity,
		RecipientThana:        body.RecipientThana,
		RecipientZip:          body.RecipientZip,
		RecipientArea:         body.RecipientArea,
		RecipientAddress:      body.RecipientAddress,
		PackageCode:           body.PackageCode,
		ParcelType:            body.ParcelType,
		RequestedDeliveryTime: body.RequestedDeliveryTime,
		PickAddress:           body.PickAddress,
		PickHub:               body.PickHub,
		Price:                 body.Price,
		NumberOfItems:         body.NumberOfItems,
		Comments:              body.Comments,
		DeliveryType:          body.DeliveryType,
		PaymentStatus:         body.PaymentStatus,
		Status: []models.OrderStatus{
			{
				ID:   primitive.NewObjectID(),
				Text: constants.Pending,
				Time: time.Now().UTC(),
			},
		},
		CreatedAt: time.Now().UTC(),
	}
	return order, nil
}

type OrderStatusUpdateReq struct {
	Text            string             `validate:"required" json:"text"`
	DeleveryBoyID   primitive.ObjectID `validate:"omitempty" json:"deleveryBoy"`
	ShopModeratorID primitive.ObjectID `validate:"omitempty" json:"shopModerator"`
	MerchantID      primitive.ObjectID `validate:"omitempty" json:"merchant"`
	AdminID         primitive.ObjectID `validate:"omitempty" json:"admin"`
	Status          string             `validate:"required" json:"status"`
}

func UpdateOrderStatus(ctx echo.Context) (*models.OrderStatus, error) {
	body := OrderStatusUpdateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}

	orderStatus := &models.OrderStatus{
		ID:              primitive.NewObjectID(),
		DeleveryBoyID:   body.DeleveryBoyID,
		ShopModeratorID: body.ShopModeratorID,
		MerchantID:      body.MerchantID,
		AdminID:         body.AdminID,
		Status:          body.Status,
		Text:            body.Text,
		Time:            time.Now().UTC(),
	}
	return orderStatus, nil
}

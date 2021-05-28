package validators

import (
	"errors"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type OrderCreateReq struct {
	RecipientName         string    `validate:"required" json:"recipientName"`
	RecipientPhone        string    `validate:"required" json:"recipientPhone"`
	RecipientCity         string    `validate:"required" json:"recipientCity"`
	RecipientThana        string    `validate:"omitempty" json:"recipientThana"`
	RecipientArea         string    `validate:"required" json:"recipientArea"`
	RecipientZip          string    `validate:"omitempty" json:"recipientZip"`
	RecipientAddress      string    `validate:"required" json:"recipientAddress"`
	PackageCode           string    `validate:"omitempty" json:"packageCode"`
	PaymentStatus         string    `validate:"required" json:"paymentStatus"`
	Price                 float64   `validate:"omitempty,number,gte=0" json:"price"`
	PercelType            string    `validate:"required" json:"percelType"`
	RequestedDeliveryTime time.Time `validate:"omitempty" json:"requestedDeliveryTime"`
	PickAddress           string    `validate:"required" json:"pickAddress"`
	PickHub               string    `validate:"required" json:"pickHub"`
	Comments              string    `validate:"omitempty,max=300" json:"comments"`
	NumberOfItems         int       `validate:"omitempty" json:"numberOfItems"`
	Weight                float32   `validate:"required,number,gt=0" json:"weight"`
	DeliveryType          string    `validate:"required" json:"deliveryType"`
}

func isPriceAcceptable(paymentStatus string, price float64) bool {
	if paymentStatus == constants.COD && price <= 0 {
		return false
	}
	return true
}

func ValidateOrderCreate(ctx echo.Context) (*models.Order, error) {
	body := OrderCreateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	if flag := isPriceAcceptable(body.PaymentStatus, body.Price); !flag {
		return nil, errors.New("price can not be zero")
	}
	created := constants.Created
	order := &models.Order{
		ID:                    primitive.NewObjectID(),
		RiderID:               nil,
		ShopModeratorID:       nil,
		MerchantID:            nil,
		RecipientName:         body.RecipientName,
		RecipientPhone:        body.RecipientPhone,
		RecipientCity:         body.RecipientCity,
		RecipientThana:        body.RecipientThana,
		RecipientZip:          body.RecipientZip,
		RecipientArea:         body.RecipientArea,
		RecipientAddress:      body.RecipientAddress,
		PackageCode:           body.PackageCode,
		PercelType:            body.PercelType,
		RequestedDeliveryTime: body.RequestedDeliveryTime,
		PickAddress:           body.PickAddress,
		PickHub:               body.PickHub,
		Price:                 body.Price,
		NumberOfItems:         body.NumberOfItems,
		Comments:              body.Comments,
		DeliveryType:          body.DeliveryType,
		PaymentStatus:         body.PaymentStatus,
		Weight:                body.Weight,
		IsCancelled:           false,
		CurrentStatus:         &created,
		IsAccepted:            false,
		Status: []models.OrderStatus{
			{
				ID:     primitive.NewObjectID(),
				Text:   "Your order have benn placed successfully",
				Status: constants.Created,
				Time:   time.Now().UTC(),
			},
		},
		CreatedAt: time.Now().UTC(),
	}
	return order, nil
}

type MultipleOrderCreateReq struct {
	Orders []OrderCreateReq `validate:"required" json:"orders"`
}

func ValidateMultipleOrderCreate(ctx echo.Context) ([]models.Order, error) {
	body := MultipleOrderCreateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	created := constants.Created
	var orders []models.Order
	for _, o := range body.Orders {
		if flag := isPriceAcceptable(o.PaymentStatus, o.Price); !flag {
			return nil, errors.New("price can not be zero")
		}
		order := models.Order{
			ID:                    primitive.NewObjectID(),
			RiderID:               nil,
			ShopModeratorID:       nil,
			MerchantID:            nil,
			RecipientName:         o.RecipientName,
			RecipientPhone:        o.RecipientPhone,
			RecipientCity:         o.RecipientCity,
			RecipientThana:        o.RecipientThana,
			RecipientZip:          o.RecipientZip,
			RecipientArea:         o.RecipientArea,
			RecipientAddress:      o.RecipientAddress,
			PackageCode:           o.PackageCode,
			PercelType:            o.PercelType,
			RequestedDeliveryTime: o.RequestedDeliveryTime,
			PickAddress:           o.PickAddress,
			PickHub:               o.PickHub,
			Price:                 o.Price,
			NumberOfItems:         o.NumberOfItems,
			Comments:              o.Comments,
			DeliveryType:          o.DeliveryType,
			PaymentStatus:         o.PaymentStatus,
			Weight:                o.Weight,
			IsCancelled:           false,
			CurrentStatus:         &created,
			IsAccepted:            false,
			Status: []models.OrderStatus{
				{
					ID:     primitive.NewObjectID(),
					Text:   "Your order have benn placed successfully",
					Status: constants.Created,
					Time:   time.Now().UTC(),
				},
			},
			CreatedAt: time.Now().UTC(),
		}
		orders = append(orders, order)
	}
	return orders, nil
}

type OrderStatusUpdateReq struct {
	Text            string             `validate:"omitempty" json:"text"`
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
		DeleveryBoyID:   &body.DeleveryBoyID,
		ShopModeratorID: &body.ShopModeratorID,
		MerchantID:      &body.MerchantID,
		AdminID:         &body.AdminID,
		Status:          body.Status,
		Text:            body.Text,
		Time:            time.Now().UTC(),
	}
	return orderStatus, nil
}

type OrderUpdateReq struct {
	RiderID               primitive.ObjectID `validate:"omitempty" json:"riderId"`
	RecipientName         string             `validate:"omitempty" json:"recipientName"`
	RecipientPhone        string             `validate:"omitempty" json:"recipientPhone"`
	RecipientCity         string             `validate:"required" json:"recipientCity"`
	RecipientThana        string             `validate:"omitempty" json:"recipientThana"`
	RecipientArea         string             `validate:"omitempty" json:"recipientArea"`
	RecipientZip          string             `validate:"omitempty" json:"recipientZip"`
	RecipientAddress      string             `validate:"omitempty" json:"recipientAddress"`
	PackageCode           string             `validate:"omitempty" json:"packageCode"`
	PaymentStatus         string             `validate:"omitempty" json:"paymentStatus"`
	Price                 float64            `validate:"omitempty,number,gte=0" json:"price"`
	PercelType            string             `validate:"omitempty" json:"percelType"`
	RequestedDeliveryTime time.Time          `validate:"omitempty" json:"requestedDeliveryTime"`
	PickAddress           string             `validate:"omitempty" json:"pickAddress"`
	PickHub               string             `validate:"omitempty" json:"pickHub"`
	Comments              string             `validate:"omitempty,max=300" json:"comments"`
	NumberOfItems         int                `validate:"omitempty" json:"numberOfItems"`
	DeliveryType          string             `validate:"required" json:"deliveryType"`
	Weight                float32            `validate:"required,number,gt=0" json:"weight"`
}

func UpdateOrder(ctx echo.Context) (*models.Order, error) {
	body := OrderUpdateReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	if flag := isPriceAcceptable(body.PaymentStatus, body.Price); !flag {
		return nil, errors.New("price can not be zero")
	}
	UserID := ctx.Get(constants.UserID).(primitive.ObjectID)
	order := &models.Order{
		RiderID:               &body.RiderID,
		RecipientName:         body.RecipientName,
		RecipientPhone:        body.RecipientPhone,
		RecipientCity:         body.RecipientCity,
		RecipientThana:        body.RecipientThana,
		RecipientZip:          body.RecipientZip,
		RecipientArea:         body.RecipientArea,
		RecipientAddress:      body.RecipientAddress,
		PackageCode:           body.PackageCode,
		PercelType:            body.PercelType,
		RequestedDeliveryTime: body.RequestedDeliveryTime,
		PickAddress:           body.PickAddress,
		PickHub:               body.PickHub,
		Price:                 body.Price,
		NumberOfItems:         body.NumberOfItems,
		Comments:              body.Comments,
		DeliveryType:          body.DeliveryType,
		PaymentStatus:         body.PaymentStatus,
		DeliveredAt:           nil,
		UpdateBy:              &UserID,
		Weight:                body.Weight,
		UpdatedAt:             time.Now().UTC(),
	}
	return order, nil
}

type OrderDeliverReq struct {
	Payment float64 `validate:"required,number,gt=-1" json:"payment"`
	Remarks string  `validate:"omitempty" json:"remarks"`
	ShopID  string  `validate:"required" json:"shopId"`
}

func ValidateOrderDeliver(ctx echo.Context) (*models.TrxHistory, error) {
	body := OrderDeliverReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	_shopID, err := primitive.ObjectIDFromHex(body.ShopID)
	if err != nil {
		return nil, err
	}
	orderID := ctx.Param("orderId")
	UserID := ctx.Get(constants.UserID).(primitive.ObjectID)
	_orderID, err := primitive.ObjectIDFromHex(orderID)
	if err != nil {
		return nil, err
	}
	trxHistory := &models.TrxHistory{
		ID:          primitive.NewObjectID(),
		Remarks:     body.Remarks,
		Payment:     body.Payment,
		CreatedBy:   UserID,
		PaymentType: models.IN,
		OrderID:     &_orderID,
		ShopID:      _shopID,
		CreatedAt:   time.Now().UTC(),
	}
	return trxHistory, nil
}

type OrderChangeReq struct {
	OrderIDs        []primitive.ObjectID `validate:"required" json:"orderIds"`
	Text            string               `validate:"omitempty" json:"text"`
	Status          string               `validate:"required" json:"status"`
	DeleveryBoyID   *primitive.ObjectID  `validate:"omitempty" json:"deleveryBoy"`
	ShopModeratorID *primitive.ObjectID  `validate:"omitempty" json:"shopModerator"`
	MerchantID      *primitive.ObjectID  `validate:"omitempty" json:"merchant"`
	AdminID         *primitive.ObjectID  `validate:"omitempty" json:"admin"`
}

func OrderChangeStatus(ctx echo.Context) (*OrderChangeReq, error) {
	body := &OrderChangeReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	userID := ctx.Get(constants.UserID).(primitive.ObjectID)
	body.AdminID = &userID
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	return body, nil
}

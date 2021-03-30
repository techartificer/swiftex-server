package validators

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RiderParcelCreate struct {
	RiderID string `json:"riderId" validate:"required"`
	OrderID string `json:"orderId" validate:"required"`
}

func ValidateRiderParcelCreate(ctx echo.Context) (*models.RiderParcel, error) {
	body := RiderParcelCreate{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	orderID, err := primitive.ObjectIDFromHex(body.OrderID)
	if err != nil {
		return nil, err
	}
	riderID, err := primitive.ObjectIDFromHex(body.RiderID)
	if err != nil {
		return nil, err
	}
	creator := ctx.Get(constants.UserID).(primitive.ObjectID)
	riderParcel := &models.RiderParcel{
		ID:         primitive.NewObjectID(),
		RiderID:    riderID,
		OrderID:    orderID,
		Status:     constants.Assigned,
		AssignedBy: creator,
		CreatedAt:  time.Now().UTC(),
	}
	return riderParcel, nil
}

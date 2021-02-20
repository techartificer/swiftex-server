package validators

import (
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants"
	"github.com/techartificer/swiftex/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type RiderCreate struct {
	Name     string `json:"name,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
	Phone    string `json:"phone,omitempty" validate:"required"`
	Contact  string `json:"contact,omitempty" validate:"required"`
	NID      string `json:"NID,omitempty" validate:"required"`
	Salary   int32  `json:"Salary,omitempty" validate:"required"`
	Address  string `json:"address,omitempty" validate:"required"`
	Remark   string `json:"remark,omitempty" validate:"required"`
}

func ValidateRiderReq(ctx echo.Context) (*models.Rider, error) {
	body := RiderCreate{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	creator := ctx.Get(constants.UserID).(primitive.ObjectID)
	rider := &models.Rider{
		ID:        primitive.NewObjectID(),
		Name:      body.Name,
		Password:  body.Password,
		Phone:     body.Phone,
		Status:    constants.Active,
		Contact:   body.Contact,
		NID:       body.NID,
		Salary:    body.Salary,
		Address:   body.Address,
		Remark:    body.Remark,
		CreatedBy: &creator,
		CreatedAt: time.Now().UTC(),
	}
	return rider, nil
}

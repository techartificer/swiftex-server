package validators

import (
	"github.com/labstack/echo/v4"
)

// ReqLogin holds login request data
type ReqLogin struct {
	Username string `json:"username,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
}

// ValidateLogin returns request body or error
func ValidateLogin(ctx echo.Context) (*ReqLogin, error) {
	body := ReqLogin{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	return &body, nil
}

package validators

import (
	"github.com/labstack/echo/v4"
)

// ReqLogin holds login request data
type ReqLogin struct {
	Phone    string `json:"phone,omitempty" validate:"required"`
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

type ForgotPasswordReq struct {
	Phone    string `json:"phone,omitempty" validate:"required"`
	Password string `json:"password,omitempty" validate:"required,min=6,max=26"`
	Token    string `json:"token,omitempty" validate:"required"`
}

// ValidateLogin returns request body or error
func ValidateForgotPassword(ctx echo.Context) (*ForgotPasswordReq, error) {
	body := ForgotPasswordReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	return &body, nil
}

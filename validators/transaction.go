package validators

import "github.com/labstack/echo/v4"

type GenerateTrxCodeReq struct {
	Amount int64 `json:"amount,omitempty"  validate:"number,gte=0"`
}

// ValidateGenerateTrxCodeReq returns request body or error
func ValidateGenerateTrxCodeReq(ctx echo.Context) (*GenerateTrxCodeReq, error) {
	body := GenerateTrxCodeReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	return &body, nil
}

type CashOutReq struct {
	TrxCode string `json:"trxCode,omitempty"  validate:"required"`
}

// ValidateCahsOutReq returns request body or error
func ValidateCahsOutReq(ctx echo.Context) (*CashOutReq, error) {
	body := CashOutReq{}
	if err := ctx.Bind(&body); err != nil {
		return nil, err
	}
	if err := GetValidationError(body); err != nil {
		return nil, err
	}
	return &body, nil
}

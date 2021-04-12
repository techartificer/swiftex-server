package validators

import "github.com/labstack/echo/v4"

type GenerateTrxCodeReq struct {
	Amount int64 `json:"amount,omitempty"  validate:"number,gte=0"`
}

// ValidateLogin returns request body or error
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

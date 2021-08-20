package response

import (
	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
)

// Response holds reponse data
type Response struct {
	Code   codes.ErrorCode `json:"code,omitempty"`
	Status int             `json:"-"`
	Title  string          `json:"title,omitempty"`
	Data   interface{}     `json:"data,omitempty"`
	Errors error           `json:"errors,omitempty"`
}

//Send after writting response send to the client
func (r *Response) Send(ctx echo.Context) error {
	ctx.Response().Header().Set("X-Platform", "SwiftEx")
	ctx.Response().Header().Set("X-Platform-Developer", "www.techartificer.com")
	ctx.Response().Header().Set("Content-Type", "application/json")

	if err := ctx.JSON(r.Status, r); err != nil {
		return err
	}
	return nil
}

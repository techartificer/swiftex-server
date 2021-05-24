package server

import (
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/techartificer/swiftex/constants/codes"
	"github.com/techartificer/swiftex/database"
	"github.com/techartificer/swiftex/lib/response"
	"github.com/ulule/limiter/v3"
)

var (
	ipRateLimiter *limiter.Limiter
)

func IPRateLimit() echo.MiddlewareFunc {
	rate := limiter.Rate{
		Period: 5 * time.Second,
		Limit:  15,
	}
	store := database.GetLimmiterStore()
	ipRateLimiter = limiter.New(*store, rate)
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			resp := response.Response{}
			ip := c.RealIP()
			limiterCtx, err := ipRateLimiter.Get(c.Request().Context(), ip)
			if err != nil {
				resp.Title = "Something went wrong"
				resp.Status = http.StatusInternalServerError
				resp.Code = codes.SomethingWentWrong
				resp.Errors = err
				return resp.Send(c)
			}

			h := c.Response().Header()
			h.Set("X-RateLimit-Limit", strconv.FormatInt(limiterCtx.Limit, 10))
			h.Set("X-RateLimit-Remaining", strconv.FormatInt(limiterCtx.Remaining, 10))
			h.Set("X-RateLimit-Reset", strconv.FormatInt(limiterCtx.Reset, 10))
			if limiterCtx.Reached {
				resp.Title = "Too Many Requests from"
				resp.Status = http.StatusTooManyRequests
				resp.Code = codes.TooManyRequest
				return resp.Send(c)
			}
			return next(c)
		}
	}
}

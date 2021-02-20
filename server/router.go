package server

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/techartificer/swiftex/api"
)

var router = echo.New()

// GetRouter returns the api router
func GetRouter() http.Handler {
	router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 6,
		Skipper: func(ctx echo.Context) bool {
			return strings.Contains(ctx.Path(), "/fs/") || strings.Contains(ctx.Path(), "/download/")
		},
	}))

	router.Pre(middleware.AddTrailingSlash())
	router.Use(middleware.Recover())
	// router.Use(echoMonitoring())
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time: ${time_rfc3339}, method: ${method}, uri: ${uri}, status: ${status}\n",
	}))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"health": "OK"})
	})

	registerV1Routes()
	return router
}

func registerV1Routes() {
	v1 := router.Group("/v1")
	auth := v1.Group("/auth")
	api.RegisterAuthRoutes(auth)
	admin := v1.Group("/admin")
	api.RegisterAdminRoutes(admin)
	merchant := v1.Group("/merchant")
	api.RegisterMerchantRoutes(merchant)
	shop := v1.Group("/shop")
	api.RegisterShopRoutes(shop)
	order := v1.Group("/order")
	api.RegisterOrderRoutes(order)
	rider := v1.Group("/rider")
	api.RegisterRiderRoutes(rider)
}

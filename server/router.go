package server

import (
	"net/http"
	_ "net/http/pprof"
	"net/url"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/techartificer/swiftex/api"
	"github.com/techartificer/swiftex/config"
)

var router = echo.New()

func wrapHandler(h http.Handler) echo.HandlerFunc {
	return func(c echo.Context) error {
		s := strings.TrimRight(c.Request().URL.String(), "/")
		c.Request().URL = &url.URL{Path: s, Host: c.Request().Host}
		c.Request().RequestURI = strings.TrimRight(c.Request().RequestURI, "/")
		h.ServeHTTP(c.Response(), c.Request())
		return nil
	}
}

// GetRouter returns the api router
func GetRouter() http.Handler {
	router.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 6,
		Skipper: func(ctx echo.Context) bool {
			return strings.Contains(ctx.Path(), "/fs/") || strings.Contains(ctx.Path(), "/download/")
		},
	}))
	router.Use(middleware.Recover())
	if config.GetServer().Env == "production" {
		router.Use(echoMonitoring())
	}
	router.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "time: ${time_rfc3339}, method: ${method}, uri: ${uri}, status: ${status}\n",
	}))
	router.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{"*"},
		AllowMethods: []string{"*"},
	}))

	router.Pre(middleware.AddTrailingSlash())
	router.GET("/", func(ctx echo.Context) error {
		return ctx.JSON(http.StatusOK, map[string]string{"health": "OK", "version": "v1.0.2"})
	})

	router.GET("/debug/pprof/*", wrapHandler(http.DefaultServeMux))

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
	trx := v1.Group("/transaction")
	api.RegisterTransactionRoutes(trx)
}

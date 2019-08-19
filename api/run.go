package api

import (
	"net/http"

	"github.com/facebookgo/grace/gracehttp"
	"github.com/go-playground/validator"
	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	AppMiddleware "github.com/ramabmtr/inventario/api/middleware"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/logger"
)

type (
	Validator struct {
		validator *validator.Validate
	}
)

func (cv *Validator) Validate(i interface{}) error {
	return cv.validator.Struct(i)
}

func Run() {
	e := echo.New()

	e.Validator = &Validator{validator: validator.New()}

	echo.NotFoundHandler = func(c echo.Context) error {
		return c.JSON(http.StatusNotFound, helper.FailResponse("Route not found"))
	}

	echo.MethodNotAllowedHandler = func(c echo.Context) error {
		return c.JSON(http.StatusMethodNotAllowed, helper.FailResponse("Method not allowed"))
	}

	e.Use(middleware.RequestIDWithConfig(middleware.RequestIDConfig{
		Generator: func() string {
			return uuid.New().String()
		},
	}))
	e.Use(middleware.CORS())
	e.Use(AppMiddleware.RequestLogger())
	e.Use(AppMiddleware.HandlePanicAndError())

	publicGroup := e.Group("")
	PublicRouter(publicGroup)

	e.Server.Addr = config.Env.App.Address

	logger.Info("API running on address ", config.Env.App.Address)

	err := gracehttp.Serve(e.Server)
	if err != nil {
		logger.WithError(err).Fatal("Fail to initialize API...")
		return
	}

	logger.Info("Server gracefully stopped... BYE")
}

package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
)

func HandlePanicAndError() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					logger.Debug("stacktrace from panic: \n" + string(debug.Stack()))
					c.JSON(http.StatusInternalServerError, helper.ErrorResponse(config.ErrDefault.Error()))
				}
			}()

			err = next(c)
			if err != nil {
				return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(config.ErrDefault.Error()))
			}

			return err
		}
	}
}

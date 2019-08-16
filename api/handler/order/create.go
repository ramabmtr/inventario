package order

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/helper"
)

func CreateOrder(c echo.Context) error {
	return c.JSON(http.StatusOK, helper.SuccessResponse())
}

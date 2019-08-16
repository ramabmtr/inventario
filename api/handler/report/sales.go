package report

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/helper"
)

func SalesReport(c echo.Context) error {
	return c.JSON(http.StatusOK, helper.SuccessResponse())
}

package report

import (
	"net/http"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

func InventoryReport(c echo.Context) error {
	var err error

	db := config.GetDatabaseClient()

	variantRepo := sqlite.NewVariantRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)

	reportSvc := service.NewReportService(variantRepo, orderRepo)

	res, err := reportSvc.GetInventoryReport(false)
	if err != nil {
		logger.WithError(err).Error("fail to process get inventory report")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.ObjectResponse(res, "inventory_report"))
}

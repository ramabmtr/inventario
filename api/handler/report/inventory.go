package report

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/repository/logger"
	"github.com/ramabmtr/inventario/service"
)

func InventoryReport(c echo.Context) error {
	var err error

	db := config.GetDatabaseClient()

	variantRepo := sqlite.NewVariantRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	trxRepo := sqlite.NewTransactionRepository(db)

	reportSvc := service.NewReportService(variantRepo, orderRepo, trxRepo)

	res, err := reportSvc.GetInventoryReport(false)
	if err != nil {
		logger.WithError(err).Error("fail to process get inventory report")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.ObjectResponse(res, "inventory_report"))
}

func InventoryReportCSV(c echo.Context) error {
	var err error

	db := config.GetDatabaseClient()

	variantRepo := sqlite.NewVariantRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	trxRepo := sqlite.NewTransactionRepository(db)

	reportSvc := service.NewReportService(variantRepo, orderRepo, trxRepo)

	res, err := reportSvc.GetInventoryReportCSV(false)
	if err != nil {
		logger.WithError(err).Error("fail to process export inventory report to csv")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=%q", "attachment", "Inventory Report.csv"))
	return c.Blob(http.StatusOK, "text/csv", res)
}

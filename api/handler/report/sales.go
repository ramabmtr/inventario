package report

import (
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

func SalesReport(c echo.Context) error {
	var err error

	startDateQuery := c.QueryParam("start_date")
	if startDateQuery == "" {
		startDateQuery = time.Now().UTC().Format(config.QueryDateFormatLayout)
	}

	endDateQuery := c.QueryParam("end_date")
	if endDateQuery == "" {
		endDateQuery = time.Now().UTC().Format(config.QueryDateFormatLayout)
	}

	startDate, err := time.Parse(config.QueryDateFormatLayout, startDateQuery)
	if err != nil {
		logger.WithError(err).Error("fail to parse start date")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	endDate, err := time.Parse(config.QueryDateFormatLayout, endDateQuery)
	if err != nil {
		logger.WithError(err).Error("fail to parse end date")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}
	// since endDate just contain the date (the time will set to default 00:00:00),
	// add time 23:59:59 to endDate (to make endDate to the end of the day).
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	db := config.GetDatabaseClient()

	variantRepo := sqlite.NewVariantRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	trxRepo := sqlite.NewTransactionRepository(db)

	reportSvc := service.NewReportService(variantRepo, orderRepo, trxRepo)

	res, err := reportSvc.GetSalesReport(&startDate, &endDate)
	if err != nil {
		logger.WithError(err).Error("fail to process get inventory report")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.ObjectResponse(res, "inventory_report"))
}

func SalesReportCSV(c echo.Context) error {
	var err error

	startDateQuery := c.QueryParam("start_date")
	if startDateQuery == "" {
		startDateQuery = time.Now().UTC().Format(config.QueryDateFormatLayout)
	}

	endDateQuery := c.QueryParam("end_date")
	if endDateQuery == "" {
		endDateQuery = time.Now().UTC().Format(config.QueryDateFormatLayout)
	}

	startDate, err := time.Parse(config.QueryDateFormatLayout, startDateQuery)
	if err != nil {
		logger.WithError(err).Error("fail to parse start date")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	endDate, err := time.Parse(config.QueryDateFormatLayout, endDateQuery)
	if err != nil {
		logger.WithError(err).Error("fail to parse end date")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}
	// since endDate just contain the date (the time will set to default 00:00:00),
	// add time 23:59:59 to endDate (to make endDate to the end of the day).
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	db := config.GetDatabaseClient()

	variantRepo := sqlite.NewVariantRepository(db)
	orderRepo := sqlite.NewOrderRepository(db)
	trxRepo := sqlite.NewTransactionRepository(db)

	reportSvc := service.NewReportService(variantRepo, orderRepo, trxRepo)

	res, err := reportSvc.GetSalesReportCSV(&startDate, &endDate)
	if err != nil {
		logger.WithError(err).Error("fail to process export inventory report to csv")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	c.Response().Header().Set("Content-Disposition", fmt.Sprintf("%s; filename=%q", "attachment", "Sales Report.csv"))
	return c.Blob(http.StatusOK, "text/csv", res)
}

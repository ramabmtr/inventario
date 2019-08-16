package order

import (
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

func GetOrderList(c echo.Context) error {
	var err error

	startDateQuery := c.QueryParam("start_date")
	if startDateQuery == "" {
		startDateQuery = time.Now().UTC().Format(config.QueryDateFormatLayout)
	}

	endDateQuery := c.QueryParam("start_date")
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

	db := config.GetDatabaseClient()

	orderRepo := sqlite.NewOrderRepository(db)

	orderSvc := service.NewOrderService(orderRepo)

	order := domain.Order{}

	orders, err := orderSvc.GetOrderList(order, &startDate, &endDate)
	if err != nil {
		logger.WithError(err).Error("fail to process get order list")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.ObjectResponse(orders, "orders"))
}

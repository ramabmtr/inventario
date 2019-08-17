package order

import (
	"errors"
	"net/http"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

func GetOrderDetail(c echo.Context) error {
	var err error

	orderID := c.Param("orderID")
	if orderID == "" {
		err = errors.New("order id is empty")
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	db := config.GetDatabaseClient()

	orderRepo := sqlite.NewOrderRepository(db)
	orderTrxRepo := sqlite.NewOrderTransactionRepository(db)
	variantRepo := sqlite.NewVariantRepository(db)

	orderSvc := service.NewOrderService(orderRepo, orderTrxRepo, variantRepo)

	order := domain.Order{
		ID: orderID,
	}

	err = orderSvc.GetOrderDetail(&order)
	if err == config.ErrNotFound {
		logger.WithError(err).Error("order not found")
		return c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
	}
	if err != nil {
		logger.WithError(err).Error("fail to process get order detail")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.ObjectResponse(order, "order"))
}

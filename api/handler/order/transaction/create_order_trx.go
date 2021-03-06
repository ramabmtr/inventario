package transaction

import (
	"errors"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

type (
	createIncomingTransactionRequestParam struct {
		ID       string `json:"id"`
		Quantity int    `json:"quantity" validate:"required"`
	}
)

func CreateIncomingTransaction(c echo.Context) error {
	var err error

	orderID := c.Param("orderID")
	if orderID == "" {
		err = errors.New("order id is empty")
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	param := new(createIncomingTransactionRequestParam)
	if err = c.Bind(param); err != nil {
		logger.WithField("validate", err.Error()).Warn("fail to bind request param")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}
	if err = c.Validate(param); err != nil {
		logger.WithField("validate", err.Error()).Warn("request param did not pas the validation")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	db := config.GetDatabaseClient()
	tx := db.Begin()
	defer func() {
		if p := recover(); p != nil {
			// A panic occurred, rollback and re-panic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// Something went wrong, rollback transaction
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit().Error
			if err != nil {
				logger.WithError(err).Error("fail to commit transaction")
				tx.Rollback()
			}
		}
	}()

	orderRepo := sqlite.NewOrderRepository(tx)
	orderTrxRepo := sqlite.NewOrderTransactionRepository(tx)
	variantRepo := sqlite.NewVariantRepository(tx)

	orderSvc := service.NewOrderService(orderRepo, orderTrxRepo, variantRepo)

	order := domain.Order{
		ID: orderID,
	}

	err = orderSvc.GetOrderDetail(&order)
	if err == config.ErrNotFound {
		logger.WithError(err).Error("order not found")
		c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
	}
	if err != nil {
		logger.WithError(err).Error("fail to process get order detail")
		c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	now := time.Now().UTC()

	orderTrx := domain.OrderTransaction{
		ID:        uuid.New().String(),
		OrderID:   orderID,
		Quantity:  param.Quantity,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	code, err := orderSvc.CreateOrderTransaction(&order, &orderTrx)
	if err != nil {
		logger.WithError(err).Error("fail to process create transaction")
		return c.JSON(code, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(orderTrx, "order_transaction"))
}

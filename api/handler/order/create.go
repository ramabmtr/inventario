package order

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/repository/logger"
	"github.com/ramabmtr/inventario/service"
)

type (
	createOrderRequestParam struct {
		ID         string  `json:"id"`
		VariantSKU string  `json:"variant_sku" validate:"required"`
		Quantity   int     `json:"quantity" validate:"required"`
		Price      float64 `json:"price" validate:"required"`
		Receipt    string  `json:"receipt"`
	}
)

func CreateOrder(c echo.Context) error {
	var err error

	param := new(createOrderRequestParam)
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

	orderID := param.ID
	if orderID == "" {
		orderID = uuid.New().String()
	}

	now := time.Now().UTC()

	order := domain.Order{
		ID:         orderID,
		VariantSKU: param.VariantSKU,
		Quantity:   param.Quantity,
		Price:      param.Price,
		Receipt:    param.Receipt,
		CreatedAt:  &now,
		UpdatedAt:  &now,
	}

	err = orderSvc.CreateOrder(&order)
	if err != nil {
		logger.WithError(err).Error("fail to process create order")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(order, "order"))
}

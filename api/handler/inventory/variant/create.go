package variant

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

type (
	CreateVariantRequestParam struct {
		SKU   string `json:"sku" validate:"required"`
		Name  string `json:"name" validate:"required_without=Size Color"`
		Size  string `json:"size" validate:"required_without=Name Color"`
		Color string `json:"color" validate:"required_without=Name Size"`
	}
)

func CreateVariant(c echo.Context) error {
	var err error

	inventoryID := c.Param("inventoryID")
	if inventoryID == "" {
		err := errors.New("inventory id is empty")
		logger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	param := new(CreateVariantRequestParam)
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

	inventoryRepo := sqlite.NewInventoryRepository(tx)
	variantRepo := sqlite.NewVariantRepository(tx)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	variant := domain.InventoryVariant{
		SKU:         param.SKU,
		InventoryID: inventoryID,
		Name:        param.Name,
		Size:        param.Size,
		Color:       param.Color,
	}

	err = inventorySvc.CreateInventoryVariant(&variant)
	if err != nil {
		logger.WithError(err).Error("fail to process create inventory")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(variant, "variant"))
}

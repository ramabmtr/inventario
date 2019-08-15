package inventory

import (
	"errors"
	"net/http"
	"time"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

type (
	updateInventoryRequestParam struct {
		Name string `json:"name" validate:"required"`
	}

	updateInventoryVariantRequestParam struct {
		Name     string `json:"name" validate:"required_without=Size Color"`
		Size     string `json:"size" validate:"required_without=Name Color"`
		Color    string `json:"color" validate:"required_without=Name Size"`
		Quantity int    `json:"quantity"`
	}
)

func UpdateInventory(c echo.Context) error {
	var err error

	inventoryID := c.Param("inventoryID")
	if inventoryID == "" {
		err := errors.New("inventory id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	param := new(updateInventoryRequestParam)
	if err := c.Bind(param); err != nil {
		config.AppLogger.WithField("validate", err.Error()).Warn("fail to bind request param")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}
	if err := c.Validate(param); err != nil {
		config.AppLogger.WithField("validate", err.Error()).Warn("request param did not pas the validation")
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
				config.AppLogger.WithError(err).Error("fail to commit transaction")
				tx.Rollback()
			}
		}
	}()

	inventoryRepo := sqlite.NewInventoryRepository(tx)
	variantRepo := sqlite.NewVariantRepository(tx)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	now := time.Now().UTC()

	i := domain.Inventory{
		ID:        inventoryID,
		Name:      param.Name,
		UpdatedAt: &now,
	}

	err = inventorySvc.UpdateInventory(&i)
	if err != nil {
		config.AppLogger.WithError(err).Error("fail to process update inventory")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(i, "inventory"))
}

func UpdateVariant(c echo.Context) error {
	var err error

	inventoryID := c.Param("inventoryID")
	if inventoryID == "" {
		err := errors.New("inventory id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	variantSKU := c.Param("variantSKU")
	if variantSKU == "" {
		err := errors.New("variant id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	param := new(updateInventoryVariantRequestParam)
	if err := c.Bind(param); err != nil {
		config.AppLogger.WithField("validate", err.Error()).Warn("fail to bind request param")
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}
	if err := c.Validate(param); err != nil {
		config.AppLogger.WithField("validate", err.Error()).Warn("request param did not pas the validation")
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
				config.AppLogger.WithError(err).Error("fail to commit transaction")
				tx.Rollback()
			}
		}
	}()

	inventoryRepo := sqlite.NewInventoryRepository(tx)
	variantRepo := sqlite.NewVariantRepository(tx)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	now := time.Now().UTC()

	variant := domain.InventoryVariant{
		SKU:         variantSKU,
		InventoryID: inventoryID,
		Name:        param.Name,
		Color:       param.Color,
		Size:        param.Size,
		Quantity:    param.Quantity,
		UpdatedAt:   &now,
	}

	err = inventorySvc.UpdateInventoryVariant(&variant)
	if err != nil {
		config.AppLogger.WithError(err).Error("fail to process update inventory")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(variant, "variant"))
}

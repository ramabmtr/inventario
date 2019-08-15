package inventory

import (
	"errors"
	"net/http"
	"strings"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

func GetInventoryDetail(c echo.Context) error {
	var err error

	inventoryID := strings.ToUpper(c.Param("inventoryID"))
	if inventoryID == "" {
		err := errors.New("inventory id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	db := config.GetDatabaseClient()

	inventoryRepo := sqlite.NewInventoryRepository(db)
	variantRepo := sqlite.NewVariantRepository(db)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	i := domain.Inventory{
		ID: inventoryID,
	}

	err = inventorySvc.GetInventoryDetail(&i)
	if err == config.ErrNotFound {
		config.AppLogger.WithError(err).Error("inventory not found")
		return c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
	}
	if err != nil {
		config.AppLogger.WithError(err).Error("fail to process get inventory detail")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(i, "inventory"))
}

func GetVariantDetail(c echo.Context) error {
	var err error

	inventoryID := strings.ToUpper(c.Param("inventoryID"))
	if inventoryID == "" {
		err := errors.New("inventory id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	variantID := strings.ToUpper(c.Param("variantID"))
	if variantID == "" {
		err := errors.New("variant id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	db := config.GetDatabaseClient()

	inventoryRepo := sqlite.NewInventoryRepository(db)
	variantRepo := sqlite.NewVariantRepository(db)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	variant := domain.Variant{
		ID:          variantID,
		InventoryID: inventoryID,
	}

	err = inventorySvc.GetInventoryVariantDetail(&variant)
	if err == config.ErrNotFound {
		config.AppLogger.WithError(err).Error("variant not found")
		return c.JSON(http.StatusNotFound, helper.FailResponse(err.Error()))
	}
	if err != nil {
		config.AppLogger.WithError(err).Error("fail to process get variant detail")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(variant, "variant"))
}

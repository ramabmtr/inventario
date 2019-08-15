package inventory

import (
	"errors"
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
	"net/http"
	"strconv"
)

func GetInventoryList(c echo.Context) error {
	var err error

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = config.DefaultLimit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = config.DefaultOffset
	}

	db := config.GetDatabaseClient()

	inventoryRepo := sqlite.NewInventoryRepository(db)
	variantRepo := sqlite.NewVariantRepository(db)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	inventories, err := inventorySvc.GetInventoryList(limit, offset)
	if err != nil {
		config.AppLogger.WithError(err).Error("fail to process get inventory list")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(inventories, "inventories"))
}

func GetVariantList(c echo.Context) error {
	var err error

	limit, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil {
		limit = config.DefaultLimit
	}

	offset, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil {
		offset = config.DefaultOffset
	}

	inventoryID := c.Param("inventoryID")
	if inventoryID == "" {
		err := errors.New("inventory id is empty")
		config.AppLogger.Error(err.Error())
		return c.JSON(http.StatusBadRequest, helper.FailResponse(err.Error()))
	}

	db := config.GetDatabaseClient()

	inventoryRepo := sqlite.NewInventoryRepository(db)
	variantRepo := sqlite.NewVariantRepository(db)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	variants, err := inventorySvc.GetVariantList(inventoryID, limit, offset)
	if err != nil {
		config.AppLogger.WithError(err).Error("fail to process get variant list")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(variants, "variants"))
}

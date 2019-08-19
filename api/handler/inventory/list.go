package inventory

import (
	"net/http"
	"strconv"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/repository/logger"
	"github.com/ramabmtr/inventario/service"
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
		logger.WithError(err).Error("fail to process get inventory list")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusOK, helper.ObjectResponse(inventories, "inventories"))
}

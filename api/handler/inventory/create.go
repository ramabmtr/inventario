package inventory

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/api/handler/inventory/variant"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/logger"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

type (
	createInventoryRequestParam struct {
		Name     string                               `json:"name" validate:"required"`
		Variants []*variant.CreateVariantRequestParam `json:"variants" validate:"required,dive,required"`
	}
)

func CreateInventory(c echo.Context) error {
	var err error

	param := new(createInventoryRequestParam)
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

	now := time.Now().UTC()

	variants := make([]domain.InventoryVariant, 0)

	inventoryID := uuid.New().String()

	for _, v := range param.Variants {
		variants = append(variants, domain.InventoryVariant{
			SKU:         v.SKU,
			InventoryID: inventoryID,
			Size:        v.Size,
			Color:       v.Color,
			CreatedAt:   &now,
			UpdatedAt:   &now,
		})
	}

	i := domain.Inventory{
		ID:        inventoryID,
		Name:      param.Name,
		CreatedAt: &now,
		UpdatedAt: &now,
		Variants:  variants,
	}

	err = inventorySvc.CreateInventory(&i)
	if err != nil {
		logger.WithError(err).Error("fail to process create inventory")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(i, "inventory"))
}

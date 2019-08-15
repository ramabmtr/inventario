package inventory

import (
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
	"github.com/ramabmtr/inventario/repository/database/sqlite"
	"github.com/ramabmtr/inventario/service"
)

type (
	createInventoryRequestParam struct {
		Name     string                 `json:"name" validate:"required"`
		Variants []*variantRequestParam `json:"variants" validate:"gt=0,required,dive,required"`
	}

	variantRequestParam struct {
		SKU      string `json:"sku" validate:"required_without=Name Size Color"`
		Name     string `json:"name" validate:"required_without=SKU Size Color"`
		Size     string `json:"size" validate:"required_without=SKU Name Color"`
		Color    string `json:"color" validate:"required_without=SKU Name Size"`
		Quantity int    `json:"quantity" validate:"required"`
	}
)

func Create(c echo.Context) error {
	var err error

	param := new(createInventoryRequestParam)
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
	variantRepo := sqlite.NewInventoryVariantRepository(tx)

	inventorySvc := service.NewInventoryService(inventoryRepo, variantRepo)

	now := time.Now().UTC()

	variants := make([]domain.InventoryVariant, 0)

	inventoryID := uuid.New().String()

	for _, v := range param.Variants {
		variants = append(variants, domain.InventoryVariant{
			ID:          uuid.New().String(),
			InventoryID: inventoryID,
			SKU:         v.SKU,
			Name:        v.Name,
			Size:        v.Size,
			Color:       v.Color,
			Quantity:    v.Quantity,
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
		config.AppLogger.WithError(err).Error("fail to process create inventory")
		return c.JSON(http.StatusInternalServerError, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(i, "inventory"))
}

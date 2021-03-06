package transaction

import (
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
	createOutgoingTransactionRequestParam struct {
		ID    string                     `json:"id"`
		Items []outgoingTransactionParam `json:"items" validate:"gt=0,dive,required"`
	}

	outgoingTransactionParam struct {
		VariantSKU string  `json:"variant_sku" validate:"required"`
		Quantity   int     `json:"quantity" validate:"required"`
		Price      float64 `json:"price" validate:"required"`
	}
)

func CreateOutgoingTransaction(c echo.Context) error {
	var err error

	param := new(createOutgoingTransactionRequestParam)
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

	trxRepo := sqlite.NewTransactionRepository(tx)
	trxItemRepo := sqlite.NewTransactionItemRepository(tx)
	variantRepo := sqlite.NewVariantRepository(tx)

	trxSvc := service.NewTransactionService(trxRepo, trxItemRepo, variantRepo)

	now := time.Now().UTC()

	trxID := param.ID
	if trxID == "" {
		trxID = uuid.New().String()
	}

	var trxItems []domain.TransactionItem

	for _, item := range param.Items {
		trxItems = append(trxItems, domain.TransactionItem{
			ID:            uuid.New().String(),
			TransactionID: trxID,
			VariantSKU:    item.VariantSKU,
			Quantity:      item.Quantity,
			Price:         item.Price,
			CreatedAt:     &now,
			UpdatedAt:     &now,
		})
	}

	trx := domain.Transaction{
		ID:        trxID,
		Items:     trxItems,
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	code, err := trxSvc.CreateOutgoingTransaction(&trx)
	if err != nil {
		logger.WithError(err).Error("fail to process create transaction")
		return c.JSON(code, helper.ErrorResponse(err.Error()))
	}

	return c.JSON(http.StatusCreated, helper.ObjectResponse(trx, "transaction"))
}

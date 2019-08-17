package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	transactionService struct {
		trx     domain.TransactionIFace
		trxItem domain.TransactionItemIFace
		variant domain.VariantIFace
	}
)

func NewTransactionService(
	trx domain.TransactionIFace,
	trxItem domain.TransactionItemIFace,
	variant domain.VariantIFace,
) *transactionService {
	return &transactionService{
		trx:     trx,
		trxItem: trxItem,
		variant: variant,
	}
}

func (c *transactionService) GetTransactionList(trx domain.Transaction, startDate, endDate *time.Time) (transactions []domain.Transaction, err error) {
	return c.trx.GetList(trx, startDate, endDate)
}

func (c *transactionService) CreateOutgoingTransaction(trx *domain.Transaction) (code int, err error) {
	err = c.trx.Create(trx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	for _, item := range trx.Items {
		variant := domain.InventoryVariant{
			SKU: item.VariantSKU,
		}

		err = c.variant.GetDetail(&variant, false)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		currentQuantity := variant.Quantity
		if item.Quantity > currentQuantity {
			return http.StatusNotAcceptable, errors.New(fmt.Sprintf("quantity for this trx exceeded the limit. max quantity allowed: %v", currentQuantity))
		}

		err = c.trxItem.Create(&item)
		if err != nil {
			return http.StatusInternalServerError, err
		}

		updatedVariant := domain.InventoryVariant{
			Quantity: currentQuantity - item.Quantity,
		}

		err = c.variant.Update(variant.SKU, &updatedVariant)
		if err != nil {
			return http.StatusInternalServerError, err
		}
	}

	return http.StatusOK, nil
}

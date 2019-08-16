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
		transaction domain.TransactionIFace
		variant     domain.VariantIFace
	}
)

func NewTransactionService(
	transaction domain.TransactionIFace,
	variant domain.VariantIFace,
) *transactionService {
	return &transactionService{
		transaction: transaction,
		variant:     variant,
	}
}

func (c *transactionService) GetTransactionList(trx domain.Transaction, startDate, endDate *time.Time) (transactions []domain.Transaction, err error) {
	return c.transaction.GetList(trx, startDate, endDate)
}

func (c *transactionService) CreateTransaction(trx *domain.Transaction) (err error) {
	return c.transaction.Create(trx)
}

func (c *transactionService) CreateIncomingTransaction(trx *domain.Transaction, order *domain.Order) (code int, err error) {
	completedQuantity := 0
	for _, v := range order.Transactions {
		completedQuantity = completedQuantity + v.Quantity
	}

	allowedQuantity := order.Quantity - completedQuantity
	if trx.Quantity > allowedQuantity {
		return http.StatusNotAcceptable, errors.New(fmt.Sprintf("quantity for this transaction exceeded the limit. max quantity allowed: %v", allowedQuantity))
	}

	err = c.transaction.Create(trx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// update variant quantity
	variant := domain.InventoryVariant{
		SKU: trx.VariantSKU,
	}
	err = c.variant.GetDetail(&variant, false)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	finalQuantity := variant.Quantity + trx.Quantity

	updatedVariant := domain.InventoryVariant{
		Quantity: finalQuantity,
	}

	err = c.variant.Update(variant.SKU, &updatedVariant)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

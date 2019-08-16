package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ramabmtr/inventario/config"
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

func (c *transactionService) CreateOutgoingTransaction(trx *domain.Transaction) (code int, err error) {
	if trx.Type != config.OutgoingTransactionType {
		return http.StatusExpectationFailed, errors.New("transaction type is not OUT")
	}

	variant := domain.InventoryVariant{
		SKU: trx.VariantSKU,
	}

	err = c.variant.GetDetail(&variant, false)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	currentQuantity := variant.Quantity
	if trx.Quantity > currentQuantity {
		return http.StatusNotAcceptable, errors.New(fmt.Sprintf("quantity for this transaction exceeded the limit. max quantity allowed: %v", currentQuantity))
	}

	err = c.transaction.Create(trx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	updatedVariant := domain.InventoryVariant{
		Quantity: currentQuantity - trx.Quantity,
	}

	err = c.variant.Update(variant.SKU, &updatedVariant)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

func (c *transactionService) CreateIncomingTransaction(trx *domain.Transaction, order *domain.Order) (code int, err error) {
	if trx.Type != config.IncomingTransactionType {
		return http.StatusExpectationFailed, errors.New("transaction type is not IN")
	}

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

	updatedVariant := domain.InventoryVariant{
		Quantity: variant.Quantity + trx.Quantity,
	}

	err = c.variant.Update(variant.SKU, &updatedVariant)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

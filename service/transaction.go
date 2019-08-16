package service

import (
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	transactionService struct {
		transaction domain.TransactionIFace
	}
)

func NewTransactionService(transaction domain.TransactionIFace) *transactionService {
	return &transactionService{
		transaction: transaction,
	}
}

func (c *transactionService) GetTransactionList(trx domain.Transaction, startDate, endDate *time.Time) (transactions []domain.Transaction, err error) {
	return c.transaction.GetList(trx, startDate, endDate)
}

func (c *transactionService) CreateTransaction(trx *domain.Transaction) (err error) {
	return c.transaction.Create(trx)
}

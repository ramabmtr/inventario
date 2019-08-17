package domain

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/ramabmtr/inventario/config"
)

type (
	Transaction struct {
		ID        string            `json:"id"`
		CreatedAt *time.Time        `json:"created_at"`
		UpdatedAt *time.Time        `json:"updated_at"`
		DeletedAt *time.Time        `json:"-"`
		Items     []TransactionItem `json:"items"`
	}

	TransactionItem struct {
		ID            string     `json:"id"`
		TransactionID string     `json:"transaction_id"`
		VariantSKU    string     `json:"variant_sku"`
		Quantity      int        `json:"quantity"`
		Price         float64    `json:"price"`
		CreatedAt     *time.Time `json:"created_at"`
		UpdatedAt     *time.Time `json:"updated_at"`
		DeletedAt     *time.Time `json:"-"`
	}

	TransactionIFace interface {
		GetList(transaction Transaction, startDate, EndDate *time.Time) (transactions []Transaction, err error)
		Create(transaction *Transaction) (err error)
	}

	TransactionItemIFace interface {
		Create(item *TransactionItem) (err error)
	}
)

func (c Transaction) MarshalJSON() ([]byte, error) {
	type Alias Transaction
	return jsoniter.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (Alias)(c),
		CreatedAt: c.CreatedAt.Format(config.ISO8601Format),
		UpdatedAt: c.UpdatedAt.Format(config.ISO8601Format),
	})
}

func (c TransactionItem) MarshalJSON() ([]byte, error) {
	type Alias TransactionItem
	return jsoniter.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	}{
		Alias:     (Alias)(c),
		CreatedAt: c.CreatedAt.Format(config.ISO8601Format),
		UpdatedAt: c.UpdatedAt.Format(config.ISO8601Format),
	})
}

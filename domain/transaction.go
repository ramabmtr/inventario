package domain

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/ramabmtr/inventario/config"
	"time"
)

type (
	Transaction struct {
		ID         string     `json:"id"`
		VariantSKU string     `json:"variant_sku"`
		OrderID    string     `json:"order_id"`
		Type       string     `json:"type"`
		Quantity   int        `json:"quantity"`
		Price      float64    `json:"price"`
		CreatedAt  *time.Time `json:"created_at"`
		UpdatedAt  *time.Time `json:"updated_at"`
		DeletedAt  *time.Time `json:"-"`
	}

	TransactionIFace interface {
		GetList(transaction Transaction, startDate, EndDate *time.Time) (transactions []Transaction, err error)
		Create(transaction *Transaction) (err error)
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

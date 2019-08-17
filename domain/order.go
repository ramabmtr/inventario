package domain

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/ramabmtr/inventario/config"
)

type (
	Order struct {
		ID           string             `json:"id"`
		VariantSKU   string             `json:"variant_sku"`
		Quantity     int                `json:"quantity"`
		Price        float64            `json:"price"`
		Receipt      string             `json:"receipt"`
		CreatedAt    *time.Time         `json:"created_at"`
		UpdatedAt    *time.Time         `json:"updated_at"`
		DeletedAt    *time.Time         `json:"-"`
		Transactions []OrderTransaction `json:"transactions"`
	}

	OrderTransaction struct {
		ID        string     `json:"id"`
		OrderID   string     `json:"order_id"`
		Quantity  int        `json:"quantity"`
		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`
	}

	OrderIFace interface {
		GetAll(order Order, fetchTransaction bool) (orders []Order, err error)
		GetDetail(order *Order) (err error)
		GetList(order Order, startDate, endDate *time.Time) (orders []Order, err error)
		Create(order *Order) (err error)
	}

	OrderTransactionIFace interface {
		Create(orderTransaction *OrderTransaction) (err error)
	}
)

func (c Order) MarshalJSON() ([]byte, error) {
	type Alias Order
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

func (c OrderTransaction) MarshalJSON() ([]byte, error) {
	type Alias OrderTransaction
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

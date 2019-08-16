package domain

import (
	jsoniter "github.com/json-iterator/go"
	"github.com/ramabmtr/inventario/config"
	"time"
)

type (
	Order struct {
		ID           string        `json:"id"`
		VariantSKU   string        `json:"variant_sku"`
		Quantity     int           `json:"quantity"`
		Price        float64       `json:"price"`
		Receipt      string        `json:"receipt"`
		CreatedAt    *time.Time    `json:"created_at"`
		UpdatedAt    *time.Time    `json:"updated_at"`
		DeletedAt    *time.Time    `json:"-"`
		Transactions []Transaction `json:"transactions"`
	}

	OrderIFace interface {
		GetAll(order Order, fetchTransaction bool) (orders []Order, err error)
		GetDetail(order *Order) (err error)
		GetList(order Order, startDate, endDate *time.Time) (orders []Order, err error)
		Create(order *Order) (err error)
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

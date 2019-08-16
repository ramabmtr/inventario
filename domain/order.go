package domain

import "time"

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
		GetList(order Order, startDate, endDate *time.Time) (orders []Order, err error)
		Create(order *Order) (err error)
	}
)

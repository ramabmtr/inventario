package domain

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/ramabmtr/inventario/config"
)

type (
	InventoryReport struct {
		CreatedAt           *time.Time            `json:"created_at"`
		TotalSKU            int                   `json:"total_sku"`
		TotalInventory      int64                 `json:"total_inventory"`
		TotalInventoryValue int64                 `json:"total_inventory_value"`
		InventoryList       []InventoryListReport `json:"inventory_list"`
	}

	InventoryListReport struct {
		SKU                  string `json:"sku"`
		Name                 string `json:"name"`
		Size                 string `json:"size"`
		Color                string `json:"color"`
		TotalAvailableItem   int    `json:"total_available_item"`
		AveragePurchasePrice int64  `json:"average_purchase_price"`
		TotalItemPrice       int64  `json:"total_item_price"`
	}
)

func (c InventoryReport) MarshalJSON() ([]byte, error) {
	type Alias InventoryReport
	return jsoniter.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     (Alias)(c),
		CreatedAt: c.CreatedAt.Format(config.ISO8601Format),
	})
}

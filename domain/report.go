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

	SalesReport struct {
		CreatedAt        *time.Time        `json:"created_at"`
		StartDate        *time.Time        `json:"start_date"`
		EndDate          *time.Time        `json:"end_date"`
		TotalSales       int64             `json:"total_sales"`
		TotalGrossProfit int64             `json:"total_gross_profit"`
		TotalSalesAmount int               `json:"total_sales_amount"`
		TotalItemSold    int               `json:"total_item_sold"`
		SalesList        []SalesListReport `json:"sales_list"`
	}

	SalesListReport struct {
		ID          string     `json:"id"`
		CreatedAt   *time.Time `json:"created_at"`
		SKU         string     `json:"sku"`
		Name        string     `json:"name"`
		Size        string     `json:"size"`
		Color       string     `json:"color"`
		Quantity    int        `json:"quantity"`
		SellPrice   int64      `json:"sell_price"`
		TotalAmount int64      `json:"total_amount"`
		BuyPrice    int64      `json:"buy_price"`
		Profit      int64      `json:"profit"`
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

func (c SalesReport) MarshalJSON() ([]byte, error) {
	type Alias SalesReport
	return jsoniter.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"`
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}{
		Alias:     (Alias)(c),
		CreatedAt: c.CreatedAt.Format(config.ISO8601Format),
		StartDate: c.StartDate.Format(config.ISO8601Format),
		EndDate:   c.EndDate.Format(config.ISO8601Format),
	})
}

func (c SalesListReport) MarshalJSON() ([]byte, error) {
	type Alias SalesListReport
	return jsoniter.Marshal(&struct {
		Alias
		CreatedAt string `json:"created_at"`
	}{
		Alias:     (Alias)(c),
		CreatedAt: c.CreatedAt.Format(config.ISO8601Format),
	})
}

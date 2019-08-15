package domain

import "time"

type (
	Inventory struct {
		ID        string     `json:"id"`
		Name      string     `json:"name"`
		CreatedAt *time.Time `json:"created_at"`
		UpdatedAt *time.Time `json:"updated_at"`
		DeletedAt *time.Time `json:"-"`
		Variants  []Variant  `json:"variants"`
	}

	Variant struct {
		ID          string     `json:"id"`
		InventoryID string     `json:"inventory_id"`
		SKU         string     `json:"sku"`
		Name        string     `json:"name"`
		Size        string     `json:"size"`
		Color       string     `json:"color"`
		Quantity    int        `json:"quantity"`
		CreatedAt   *time.Time `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at"`
		DeletedAt   *time.Time `json:"-"`
		Parent      Inventory  `json:"-"`
	}

	// you can use any database to manage inventory by implementing all interface below
	// in this app, I use sqlite3. Take a look in `./repository/database/sqlite/inventory.go`
	InventoryIFace interface {
		GetList(limit, offset int, fetchVariant bool) (inventories []Inventory, err error)
		GetDetail(inventory *Inventory, fetchVariant bool) (err error)
		Create(inventory *Inventory) (err error)
		Update(inventory *Inventory) (err error)
	}

	VariantIFace interface {
		GetList(inventoryID string, limit, offset int) (variants []Variant, err error)
		GetDetail(inventory *Variant, fetchParent bool) (err error)
		Create(variant *Variant) (err error)
		Update(variant *Variant) (err error)
	}
)

package domain

import (
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/ramabmtr/inventario/config"
)

type (
	Inventory struct {
		ID        string             `json:"id"`
		Name      string             `json:"name"`
		CreatedAt *time.Time         `json:"created_at"`
		UpdatedAt *time.Time         `json:"updated_at"`
		DeletedAt *time.Time         `json:"-"`
		Variants  []InventoryVariant `json:"variants"`
	}

	InventoryVariant struct {
		SKU         string     `json:"sku"`
		InventoryID string     `json:"inventory_id"`
		Name        string     `json:"name"`
		Size        string     `json:"size"`
		Color       string     `json:"color"`
		Quantity    int        `json:"quantity"`
		CreatedAt   *time.Time `json:"created_at"`
		UpdatedAt   *time.Time `json:"updated_at"`
		DeletedAt   *time.Time `json:"-"`
		Inventory   *Inventory `json:"inventory,omitempty"`
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
		GetAll(variant InventoryVariant, showEmptyStock bool) (variants []InventoryVariant, err error)
		GetList(variant InventoryVariant, limit, offset int) (variants []InventoryVariant, err error)
		GetDetail(variant *InventoryVariant, fetchInventory bool) (err error)
		Create(variant *InventoryVariant) (err error)
		Update(sku string, variant *InventoryVariant) (err error)
	}
)

func (c Inventory) MarshalJSON() ([]byte, error) {
	type Alias Inventory
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

func (c InventoryVariant) MarshalJSON() ([]byte, error) {
	type Alias InventoryVariant
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

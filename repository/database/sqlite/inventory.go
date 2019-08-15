package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/ramabmtr/inventario/domain"
)

type (
	inventoryRepository struct {
		db *gorm.DB
	}

	inventoryVariantRepository struct {
		db *gorm.DB
	}
)

// NewInventoryRepository implements domain.InventoryIFace
// to manage the inventory data with sqlite3 database
func NewInventoryRepository(db *gorm.DB) domain.InventoryIFace {
	return &inventoryRepository{
		db: db,
	}
}

func (c *inventoryRepository) GetList(limit, offset int, fetchVariant bool) (inventories []domain.Inventory, err error) {
	return
}

func (c *inventoryRepository) GetDetail(inventory *domain.Inventory, fetchVariant bool) (err error) {
	return
}

func (c *inventoryRepository) Create(inventory *domain.Inventory) (err error) {
	return
}

func (c *inventoryRepository) Update(inventory *domain.Inventory) (err error) {
	return
}

// NewInventoryVariantRepository implements domain.InventoryVariantIFace
// to manage the inventory variant data with sqlite3 database
func NewInventoryVariantRepository(db *gorm.DB) domain.InventoryVariantIFace {
	return &inventoryVariantRepository{
		db: db,
	}
}

func (c *inventoryVariantRepository) GetDetail(variant *domain.InventoryVariant, fetchParent bool) (err error) {
	return
}

func (c *inventoryVariantRepository) Create(variant *domain.InventoryVariant) (err error) {
	return
}

func (c *inventoryVariantRepository) Update(variant *domain.InventoryVariant) (err error) {
	return
}

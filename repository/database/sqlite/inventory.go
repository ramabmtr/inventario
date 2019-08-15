package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/ramabmtr/inventario/domain"
)

type (
	inventoryRepository struct {
		db *gorm.DB
	}

	variantRepository struct {
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
	// take out the variant from the struct because the behaviour of lib
	// that will automatically insert the association struct
	variant := inventory.Variants
	inventory.Variants = nil
	defer func() {
		inventory.Variants = variant
	}()
	return c.db.Create(inventory).Error
}

func (c *inventoryRepository) Update(inventory *domain.Inventory) (err error) {
	return
}

// NewVariantRepository implements domain.VariantIFace
// to manage the inventory variant data with sqlite3 database
func NewVariantRepository(db *gorm.DB) domain.VariantIFace {
	return &variantRepository{
		db: db,
	}
}

func (c *variantRepository) GetList(inventoryID string, limit, offset int) (variants []domain.Variant, err error) {
	return
}

func (c *variantRepository) GetDetail(variant *domain.Variant, fetchParent bool) (err error) {
	return
}

func (c *variantRepository) Create(variant *domain.Variant) (err error) {
	return c.db.Create(variant).Error
}

func (c *variantRepository) Update(variant *domain.Variant) (err error) {
	return
}

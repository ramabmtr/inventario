package sqlite

import (
	"github.com/jinzhu/gorm"
	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/helper"
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
	inventories = make([]domain.Inventory, 0)
	q := c.db
	if fetchVariant {
		q = q.Preload("Variants")
	}
	err = q.Find(&inventories).Limit(limit).Offset(offset).Error
	return inventories, helper.TranslateSqliteError(err)
}

func (c *inventoryRepository) GetDetail(inventory *domain.Inventory, fetchVariant bool) (err error) {
	q := c.db
	if fetchVariant {
		q = q.Preload("Variants")
	}
	err = q.Where(inventory).First(inventory).Error
	return helper.TranslateSqliteError(err)
}

func (c *inventoryRepository) Create(inventory *domain.Inventory) (err error) {
	err = c.db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Create(inventory).Error
	return helper.TranslateSqliteError(err)
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

func (c *variantRepository) GetAll(variant domain.InventoryVariant, showEmptyStock bool) (variants []domain.InventoryVariant, err error) {
	variants = make([]domain.InventoryVariant, 0)
	q := c.db.Preload("Inventory").
		Where(&variant)
	if !showEmptyStock {
		q = q.Where("quantity > ?", 0)
	}
	err = q.Find(&variants).Error
	return variants, helper.TranslateSqliteError(err)
}

func (c *variantRepository) GetList(variant domain.InventoryVariant, limit, offset int) (variants []domain.InventoryVariant, err error) {
	variants = make([]domain.InventoryVariant, 0)
	err = c.db.Where(&variant).Find(&variants).Limit(limit).Offset(offset).Error
	return variants, helper.TranslateSqliteError(err)
}

func (c *variantRepository) GetDetail(variant *domain.InventoryVariant, fetchInventory bool) (err error) {
	q := c.db
	if fetchInventory {
		q = q.Preload("Inventory")
	}
	err = q.Where(variant).First(variant).Error
	return helper.TranslateSqliteError(err)
}

func (c *variantRepository) Create(variant *domain.InventoryVariant) (err error) {
	return c.db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Create(variant).Error
}

func (c *variantRepository) Update(sku string, variant *domain.InventoryVariant) (err error) {
	v := domain.InventoryVariant{
		SKU: sku,
	}

	return c.db.Model(v).Where(v).Update(variant).Error
}

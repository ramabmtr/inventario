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

func (c *variantRepository) GetList(inventoryID string, limit, offset int) (variants []domain.InventoryVariant, err error) {
	variants = make([]domain.InventoryVariant, 0)
	queryVariant := domain.InventoryVariant{
		InventoryID: inventoryID,
	}
	err = c.db.Where(&queryVariant).Find(&variants).Limit(limit).Offset(offset).Error
	return variants, helper.TranslateSqliteError(err)
}

func (c *variantRepository) GetDetail(variant *domain.InventoryVariant, fetchParent bool) (err error) {
	err = c.db.Where(variant).First(variant).Error
	if err != nil {
		return helper.TranslateSqliteError(err)
	}
	i := domain.Inventory{
		ID: variant.InventoryID,
	}
	err = c.db.Where(&i).First(&i).Error
	if err != nil {
		return helper.TranslateSqliteError(err)
	}

	variant.Parent = &i

	return nil
}

func (c *variantRepository) Create(variant *domain.InventoryVariant) (err error) {
	return c.db.Create(variant).Error
}

func (c *variantRepository) Update(sku string, variant *domain.InventoryVariant) (err error) {
	v := domain.InventoryVariant{
		SKU: sku,
	}

	return c.db.Model(v).Where(v).Update(variant).Error
}

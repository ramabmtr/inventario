package service

import (
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	inventoryService struct {
		inventory domain.InventoryIFace
		variant   domain.VariantIFace
	}
)

func NewInventoryService(inventory domain.InventoryIFace, variant domain.VariantIFace) *inventoryService {
	return &inventoryService{
		inventory: inventory,
		variant:   variant,
	}
}

func (c *inventoryService) GetInventoryList(limit, offset int) (inventories []domain.Inventory, err error) {
	return c.inventory.GetList(limit, offset, true)
}

func (c *inventoryService) GetInventoryDetail(inventory *domain.Inventory) (err error) {
	return c.inventory.GetDetail(inventory, true)
}

func (c *inventoryService) CreateInventory(inventory *domain.Inventory) (err error) {
	if err = c.inventory.Create(inventory); err != nil {
		return err
	}

	// insert the variants if any
	for _, variant := range inventory.Variants {
		if err = c.variant.Create(&variant); err != nil {
			return err
		}
	}
	return
}

func (c *inventoryService) UpdateInventory(inventory *domain.Inventory) (err error) {
	now := time.Now().UTC()
	inventory.UpdatedAt = &now
	return c.inventory.Update(inventory)
}

func (c *inventoryService) GetVariantList(inventoryID string, limit, offset int) (variants []domain.Variant, err error) {
	return c.variant.GetList(inventoryID, limit, offset)
}

func (c *inventoryService) GetInventoryVariantDetail(variant *domain.Variant) (err error) {
	return c.variant.GetDetail(variant, true)
}

func (c *inventoryService) CreateInventoryVariant(variant *domain.Variant) (err error) {
	return c.variant.Create(variant)
}

func (c *inventoryService) UpdateInventoryVariant(variant *domain.Variant) (err error) {
	now := time.Now().UTC()
	variant.UpdatedAt = &now
	return c.variant.Update(variant)
}

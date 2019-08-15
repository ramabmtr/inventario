package service

import (
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	inventoryService struct {
		inventory domain.InventoryIFace
		variant   domain.InventoryVariantIFace
	}
)

func NewInventoryService(inventory domain.InventoryIFace, variant domain.InventoryVariantIFace) *inventoryService {
	return &inventoryService{
		inventory: inventory,
		variant:   variant,
	}
}

func (c *inventoryService) Get(inventory *domain.Inventory) (err error) {
	return c.inventory.GetDetail(inventory, true)
}

func (c *inventoryService) GetAllInventory(limit, offset int) (inventories []domain.Inventory, err error) {
	return c.inventory.GetList(limit, offset, true)
}

func (c *inventoryService) CreateInventory(inventory *domain.Inventory) (err error) {
	if err = c.inventory.Create(inventory); err != nil {
		return err
	}

	// insert the variants
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

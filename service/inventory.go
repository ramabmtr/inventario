package service

import (
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	inventoryService struct {
		inventory domain.InventoryIFace
	}
)

func NewInventoryService(inventory domain.InventoryIFace) *inventoryService {
	return &inventoryService{
		inventory: inventory,
	}
}

func (c *inventoryService) Get(inventory *domain.Inventory) (err error) {
	return c.inventory.GetDetail(inventory, true)
}

func (c *inventoryService) GetAllInventory(limit, offset int) (inventories []domain.Inventory, err error) {
	return c.inventory.GetList(limit, offset, true)
}

func (c *inventoryService) CreateInventory(inventory *domain.Inventory) (err error) {
	return c.inventory.Create(inventory)
}

func (c *inventoryService) UpdateInventory(inventory *domain.Inventory) (err error) {
	now := time.Now().UTC()
	inventory.UpdatedAt = &now
	return c.inventory.Update(inventory)
}

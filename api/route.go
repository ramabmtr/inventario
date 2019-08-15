package api

import (
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/api/handler/healthcheck"
	"github.com/ramabmtr/inventario/api/handler/inventory"
)

func PublicRouter(g *echo.Group) {
	g.GET(
		"/ping",
		healthcheck.Ping,
	)

	g.GET(
		"/inventories",
		inventory.GetInventoryList,
	)

	g.POST(
		"/inventories",
		inventory.CreateInventory,
	)

	g.GET(
		"/inventories/:inventoryID",
		inventory.GetInventoryDetail,
	)

	g.PUT(
		"/inventories/:inventoryID",
		inventory.UpdateInventory,
	)

	g.GET(
		"/inventories/:inventoryID/variants",
		inventory.GetVariantList,
	)

	g.POST(
		"/inventories/:inventoryID/variants",
		inventory.CreateVariant,
	)

	g.GET(
		"/inventories/:inventoryID/variants/:variantSKU",
		inventory.GetVariantDetail,
	)

	g.PUT(
		"/inventories/:inventoryID/variants/:variantID",
		inventory.UpdateVariant,
	)
}

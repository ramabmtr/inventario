package api

import (
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/api/handler/healthcheck"
	"github.com/ramabmtr/inventario/api/handler/inventory"
	"github.com/ramabmtr/inventario/api/handler/order"
	"github.com/ramabmtr/inventario/api/handler/report"
	"github.com/ramabmtr/inventario/api/handler/transaction"
)

func PublicRouter(g *echo.Group) {
	g.GET("/ping", healthcheck.Ping)

	g.GET("/inventories", inventory.GetInventoryList)
	g.POST("/inventories", inventory.CreateInventory)
	g.GET("/inventories/:inventoryID", inventory.GetInventoryDetail)
	g.PUT("/inventories/:inventoryID", inventory.UpdateInventory)
	g.GET("/inventories/:inventoryID/variants", inventory.GetVariantList)
	g.POST("/inventories/:inventoryID/variants", inventory.CreateVariant)
	g.GET("/inventories/:inventoryID/variants/:variantSKU", inventory.GetVariantDetail)
	g.PUT("/inventories/:inventoryID/variants/:variantID", inventory.UpdateVariant)

	g.GET("/orders", order.GetOrderList)
	g.POST("/orders", order.CreateOrder)

	g.GET("/transactions", transaction.GetTransactionList)
	g.POST("/transactions", transaction.CreateTransaction)

	g.GET("/reports/inventory", report.InventoryReport)
	g.GET("/reports/sales", report.SalesReport)
}

package api

import (
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/api/handler/healthcheck"
	"github.com/ramabmtr/inventario/api/handler/inventory"
	"github.com/ramabmtr/inventario/api/handler/inventory/variant"
	"github.com/ramabmtr/inventario/api/handler/order"
	transactionOrder "github.com/ramabmtr/inventario/api/handler/order/transaction"
	"github.com/ramabmtr/inventario/api/handler/report"
	"github.com/ramabmtr/inventario/api/handler/transaction"
)

func PublicRouter(g *echo.Group) {
	g.GET("/ping", healthcheck.Ping)

	g.GET("/inventories", inventory.GetInventoryList)
	g.POST("/inventories", inventory.CreateInventory)
	g.GET("/inventories/:inventoryID", inventory.GetInventoryDetail)
	g.PUT("/inventories/:inventoryID", inventory.UpdateInventory)

	g.GET("/inventories/:inventoryID/variants", variant.GetVariantList)
	g.POST("/inventories/:inventoryID/variants", variant.CreateVariant)
	g.GET("/inventories/:inventoryID/variants/:variantSKU", variant.GetVariantDetail)
	g.PUT("/inventories/:inventoryID/variants/:variantID", variant.UpdateVariant)

	g.GET("/orders", order.GetOrderList)
	g.POST("/orders", order.CreateOrder)
	g.GET("/orders/:orderID", order.GetOrderDetail)
	g.POST("/orders/:orderID/transactions", transactionOrder.CreateIncomingTransaction)

	g.GET("/transactions", transaction.GetTransactionList)
	g.POST("/transactions", transaction.CreateOutgoingTransaction)

	g.GET("/reports/inventories", report.InventoryReport)
	g.GET("/reports/sales", report.SalesReport)

	g.GET("/reports/inventories/exportcsv", report.InventoryReportCSV)
	g.GET("/reports/sales/exportcsv", report.SalesReportCSV)
}

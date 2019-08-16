package service

import (
	"time"

	"github.com/ramabmtr/inventario/domain"
	"github.com/ramabmtr/inventario/logger"
	"github.com/shopspring/decimal"
)

type (
	reportService struct {
		variant domain.VariantIFace
		order   domain.OrderIFace
	}
)

func NewReportService(variant domain.VariantIFace, order domain.OrderIFace) *reportService {
	return &reportService{
		variant: variant,
		order:   order,
	}
}

func (c *reportService) GetInventoryReport(showEmptyStock bool) (res *domain.InventoryReport, err error) {
	variant := domain.InventoryVariant{}
	variants, err := c.variant.GetAll(variant, showEmptyStock)
	if err != nil {
		return res, err
	}

	inventoryList := make([]domain.InventoryListReport, 0)
	totalInventory := decimal.NewFromFloat(0)
	totalInventoryPrice := decimal.NewFromFloat(0)

	for _, v := range variants {
		logger.Debug(v)
		order := domain.Order{
			VariantSKU: v.SKU,
		}

		orders, err := c.order.GetAll(order, false)
		if err != nil {
			return res, err
		}

		totalOrderQuantity := decimal.NewFromFloat(0)
		totalOrderPrice := decimal.NewFromFloat(0)

		for _, o := range orders {
			q := decimal.NewFromFloat(float64(o.Quantity))
			p := decimal.NewFromFloat(o.Price)
			tp := p.Mul(q)

			totalOrderQuantity = totalOrderQuantity.Add(q)
			totalOrderPrice = totalOrderPrice.Add(tp)
		}

		averagePrice := totalOrderPrice.Div(totalOrderQuantity)
		variantQuantity := decimal.NewFromFloat(float64(v.Quantity))
		totalItemPrice := variantQuantity.Mul(averagePrice)

		iList := domain.InventoryListReport{
			SKU:                  v.SKU,
			Name:                 v.Inventory.Name,
			Size:                 v.Size,
			Color:                v.Color,
			TotalAvailableItem:   v.Quantity,
			AveragePurchasePrice: averagePrice.Ceil().IntPart(),
			TotalItemPrice:       totalItemPrice.Ceil().IntPart(),
		}

		inventoryList = append(inventoryList, iList)
		totalInventory = totalInventory.Add(variantQuantity)
		totalInventoryPrice = totalInventoryPrice.Add(totalItemPrice)
	}

	now := time.Now().UTC()

	res = &domain.InventoryReport{
		CreatedAt:           &now,
		TotalSKU:            len(variants),
		TotalInventory:      totalInventory.IntPart(),
		TotalInventoryValue: totalInventoryPrice.Ceil().IntPart(),
		InventoryList:       inventoryList,
	}

	return
}

func (c *reportService) GetOrderList(order domain.Order, startDate, endDate *time.Time) (orders []domain.Order, err error) {
	return c.order.GetList(order, startDate, endDate)
}

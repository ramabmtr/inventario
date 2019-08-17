package service

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"time"

	jsoniter "github.com/json-iterator/go"
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

// data:
//     Data tobe written in csv. If header is nil, data will written in csv in the same order as the map data
// header:
//     Add header data on csv, you can set position or hide some parts in data
// autoDetectHeader:
//     Will automatically set header based on key on data. this field ignored when header is not nil
func (c *reportService) writeCSV(data []map[string]interface{}, header []string, autoDetectHeader bool) (res []byte, err error) {
	var records [][]string

	for k, row := range data {
		var content []string

		switch {
		case header != nil:
			for _, val := range header {
				content = append(content, fmt.Sprintf("%v", row[val]))
			}
		case autoDetectHeader:
			fallthrough
		default:
			for keyName, val := range row {
				if header == nil && autoDetectHeader {
					header = append(header, keyName)
				}
				content = append(content, fmt.Sprintf("%v", val))
			}
		}

		// write header
		if header != nil && k == 0 {
			records = append(records, header)
		}
		records = append(records, content)
	}

	buffer := new(bytes.Buffer)
	csvWriter := csv.NewWriter(buffer)
	err = csvWriter.WriteAll(records)
	return buffer.Bytes(), err
}

func (c *reportService) GetInventoryReportCSV(showEmptyStock bool) (res []byte, err error) {
	r, err := c.GetInventoryReport(showEmptyStock)
	if err != nil {
		return
	}

	var records []map[string]interface{}
	recordsByte, _ := jsoniter.Marshal(r.InventoryList)
	err = jsoniter.Unmarshal(recordsByte, &records)
	if err != nil {
		return res, err
	}

	header := []string{"sku", "name", "size", "color", "total_available_item", "average_purchase_price", "total_item_price"}

	res, err = c.writeCSV(records, header, false)
	return
}

func (c *reportService) GetOrderList(order domain.Order, startDate, endDate *time.Time) (orders []domain.Order, err error) {
	return c.order.GetList(order, startDate, endDate)
}

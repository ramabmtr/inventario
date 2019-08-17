package service

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	orderService struct {
		order    domain.OrderIFace
		orderTrx domain.OrderTransactionIFace
		variant  domain.VariantIFace
	}
)

func NewOrderService(order domain.OrderIFace, orderTrx domain.OrderTransactionIFace, variant domain.VariantIFace) *orderService {
	return &orderService{
		order:    order,
		orderTrx: orderTrx,
		variant:  variant,
	}
}

func (c *orderService) GetOrderDetail(order *domain.Order) (err error) {
	return c.order.GetDetail(order)
}

func (c *orderService) GetOrderList(order domain.Order, startDate, endDate *time.Time) (orders []domain.Order, err error) {
	return c.order.GetList(order, startDate, endDate)
}

func (c *orderService) CreateOrder(order *domain.Order) (err error) {
	return c.order.Create(order)
}

func (c *orderService) CreateOrderTransaction(order *domain.Order, orderTrx *domain.OrderTransaction) (code int, err error) {
	completedQuantity := 0
	for _, v := range order.Transactions {
		completedQuantity = completedQuantity + v.Quantity
	}

	allowedQuantity := order.Quantity - completedQuantity
	if orderTrx.Quantity > allowedQuantity {
		return http.StatusNotAcceptable, errors.New(fmt.Sprintf("quantity for this trx exceeded the limit. max quantity allowed: %v", allowedQuantity))
	}

	err = c.orderTrx.Create(orderTrx)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	// update variant quantity
	variant := domain.InventoryVariant{
		SKU: order.VariantSKU,
	}
	err = c.variant.GetDetail(&variant, false)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	updatedVariant := domain.InventoryVariant{
		Quantity: variant.Quantity + orderTrx.Quantity,
	}

	err = c.variant.Update(variant.SKU, &updatedVariant)
	if err != nil {
		return http.StatusInternalServerError, err
	}

	return http.StatusOK, nil
}

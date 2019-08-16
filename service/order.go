package service

import (
	"time"

	"github.com/ramabmtr/inventario/domain"
)

type (
	orderService struct {
		order domain.OrderIFace
	}
)

func NewOrderService(order domain.OrderIFace) *orderService {
	return &orderService{
		order: order,
	}
}

func (c *orderService) GetOrderList(order domain.Order, startDate, endDate *time.Time) (orders []domain.Order, err error) {
	return c.order.GetList(order, startDate, endDate)
}

func (c *orderService) CreateOrder(order *domain.Order) (err error) {
	return c.order.Create(order)
}

package sqlite

import (
	"time"

	"github.com/ramabmtr/inventario/helper"

	"github.com/jinzhu/gorm"
	"github.com/ramabmtr/inventario/domain"
)

type (
	orderRepository struct {
		db *gorm.DB
	}
)

// NewOrderRepository implements domain.OrderIFace
// to manage the order data with sqlite3 database
func NewOrderRepository(db *gorm.DB) domain.OrderIFace {
	return &orderRepository{
		db: db,
	}
}

func (c *orderRepository) GetDetail(order *domain.Order) (err error) {
	err = c.db.Where(order).
		Preload("Transactions").
		First(&order).Error

	return helper.TranslateSqliteError(err)
}

func (c *orderRepository) GetList(order domain.Order, startDate, endDate *time.Time) (orders []domain.Order, err error) {
	err = c.db.Where(order).
		Where("created_at >= ?", startDate).
		Where("created_at <= ?", endDate).
		Order("created_at DESC").
		Preload("Transactions").
		Find(&orders).Error
	return orders, err
}

func (c *orderRepository) Create(order *domain.Order) (err error) {
	return c.db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Create(order).Error
}

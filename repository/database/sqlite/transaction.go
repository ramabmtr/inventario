package sqlite

import (
	"time"

	"github.com/jinzhu/gorm"
	"github.com/ramabmtr/inventario/domain"
)

type (
	transactionRepository struct {
		db *gorm.DB
	}

	transactionItemRepository struct {
		db *gorm.DB
	}
)

// NewOrderRepository implements domain.TransactionIFace
// to manage the transaction data with sqlite3 database
func NewTransactionRepository(db *gorm.DB) domain.TransactionIFace {
	return &transactionRepository{
		db: db,
	}
}

func (c *transactionRepository) GetList(trx domain.Transaction, startDate, endDate *time.Time) (transactions []domain.Transaction, err error) {
	err = c.db.Where(trx).
		Preload("Items").
		Where("created_at >= ?", startDate).
		Where("created_at <= ?", endDate).
		Order("created_at DESC").
		Find(&transactions).Error
	return transactions, err
}

func (c *transactionRepository) Create(trx *domain.Transaction) (err error) {
	return c.db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Create(trx).Error
}

func NewTransactionItemRepository(db *gorm.DB) domain.TransactionItemIFace {
	return &transactionItemRepository{
		db: db,
	}
}

func (c *transactionItemRepository) Create(trxItem *domain.TransactionItem) (err error) {
	return c.db.Set("gorm:association_autoupdate", false).
		Set("gorm:association_autocreate", false).
		Create(trxItem).Error
}

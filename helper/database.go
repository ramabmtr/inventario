package helper

import (
	"github.com/jinzhu/gorm"
	"github.com/ramabmtr/inventario/config"
)

func TranslateSqliteError(err error) error {
	if err == gorm.ErrRecordNotFound {
		return config.ErrNotFound
	}
	return err
}

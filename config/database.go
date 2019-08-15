package config

import (
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

var (
	dbOnce sync.Once
	db     *gorm.DB
)

func InitDatabaseClient() (*gorm.DB, error) {
	var err error
	dbOnce.Do(func() {
		db, err = gorm.Open(Env.Database.Engine, Env.Database.URL)
		if err != nil {
			return
		}
		db.LogMode(Env.App.Debug)
		if Env.Database.Engine == DatabaseEngineSqlite3 {
			// the 3rd party lib open the DB with `SQLITE_OPEN_FULLMUTEX` flag
			// set max open connection to 1 for Sqlite because of this behaviour
			// to avoid random errors that tell that the database is locked.
			db.DB().SetMaxOpenConns(1)
		}
	})

	err = db.DB().Ping()
	return db, err
}

func GetDatabaseClient() *gorm.DB {
	if db == nil {
		db, _ = InitDatabaseClient()
	}
	return db
}

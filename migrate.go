// +build migrate

package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/repository/logger"
)

const (
	migrationFilePath    = "./migration"
	migrationSchemaTable = "migration_schema"
)

var (
	once sync.Once
	db   *gorm.DB
)

func getConn() *gorm.DB {
	var err error

	once.Do(func() {
		db, err = gorm.Open(config.Env.Database.Engine, config.Env.Database.URL)
		if err != nil {
			logger.Fatal("fail to open connection. Err: ", err.Error())
		}
	})
	return db
}

func init() {
	config.InitEnvVar()
	config.InitLogger()
}

func main() {
	db := getConn()
	defer db.Close()

	args := os.Args[1:]

	if len(args) > 0 {
		switch strings.ToLower(args[0]) {
		case "create":
			if len(args) < 2 {
				logger.Warn("migration filename is not provided")
				logger.Warn("usage: migrate create \"migration name\"")
				return
			}

			filename := toSnakeCase(args[1])

			timeStamp := time.Now().Unix()
			upMigrationFilePath := fmt.Sprintf("%s/%d_%s.up.sql", migrationFilePath, timeStamp, filename)
			downMigrationFilePath := fmt.Sprintf("%s/%d_%s.down.sql", migrationFilePath, timeStamp, filename)

			if err := createFile(upMigrationFilePath); err != nil {
				logger.Fatal("error create migration up file")
			}

			if err := createFile(downMigrationFilePath); err != nil {
				_ = os.Remove(upMigrationFilePath)
				logger.Fatal("error create migration down file")
			}

		case "up":
			files, err := ioutil.ReadDir(migrationFilePath)
			if err != nil {
				logger.Fatal("cannot read migration dir. Err: ", err.Error())
			}

			fileMap := make(map[string]string)
			fileSort := make([]string, 0)

			for _, file := range files {
				key := strings.Split(file.Name(), "_")[0]
				ops := strings.Split(file.Name(), ".")[1]
				if ops != "up" {
					continue
				}
				fileMap[key] = file.Name()
				fileSort = append(fileSort, key)
			}

			sort.Strings(fileSort)

			// get latest migration from db
			var migrationSchema string
			err = db.Raw("SELECT name FROM sqlite_master WHERE type='table' AND name=$1", migrationSchemaTable).Row().Scan(&migrationSchema)
			if err != nil && err != sql.ErrNoRows {
				logger.Fatal("fail to check migration_schema table is exist or not. Err: ", err.Error())
			}

			if migrationSchema == "" {
				logger.Info("table migration_schema is missing, creating...")
				err = db.Exec(fmt.Sprintf("CREATE TABLE %s (schema VARCHAR NOT NULL)", migrationSchemaTable)).Error
				if err != nil {
					logger.Fatal("cannot create migration_schema table. Err: ", err.Error())
				}

				err = db.Exec(fmt.Sprintf("INSERT INTO %s (schema) VALUES ('')", migrationSchemaTable)).Error
				if err != nil {
					logger.Fatal("cannot insert to migration_schema table. Err: ", err.Error())
				}
			}

			var latestMigration string
			err = db.Raw(fmt.Sprintf("SELECT schema FROM %s LIMIT 1", migrationSchemaTable)).Row().Scan(&latestMigration)
			if err != nil && err != sql.ErrNoRows {
				logger.Fatal("fail to check latest migration in migration_schema table. Err: ", err.Error())
			}

			runMigration := false

			for _, migrationVer := range fileSort {
				if latestMigration == "" {
					runMigration = true
				}

				if runMigration {
					err = runMigrationFile(migrationVer, fmt.Sprintf("%s/%s", migrationFilePath, fileMap[migrationVer]))
					if err != nil {
						logger.Fatal("fail to run migration, Err: ", err.Error())
					}
				}

				if latestMigration == migrationVer {
					runMigration = true
				}
			}

			logger.Info("Done!")
		}
	} else {
		logger.Warn("usage: migrate [create] [up] [down]")
	}
}

func createFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}

func toSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.Replace(snake, " ", "_", -1)
	return strings.ToLower(snake)
}

func runMigrationFile(ver, filepath string) error {
	db := getConn()
	tx := db.Begin()
	var err error

	logger.Info("running migration ver: ", ver)

	defer func() {
		if p := recover(); p != nil {
			// A panic occurred, rollback and re-panic
			tx.Rollback()
			panic(p)
		} else if err != nil {
			// Something went wrong, rollback transaction
			tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit().Error
			if err != nil {
				logger.Error("fail to commit transaction. Err: ", err.Error())
				tx.Rollback()
			}
		}
	}()

	sqlQuery, err := ioutil.ReadFile(filepath)
	if err != nil {
		logger.Error("file reading error, Err: ", err.Error())
		return err
	}

	if err = tx.Exec(string(sqlQuery)).Error; err != nil {
		logger.Error("error running query. Err: ", err.Error())
		return err
	}

	// update latest version to schema
	if err = tx.Exec(fmt.Sprintf("UPDATE %s SET schema = ?", migrationSchemaTable), ver).Error; err != nil {
		logger.Error("error running query. Err: ", err.Error())
		return err
	}

	return err
}

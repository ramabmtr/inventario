package main

import (
	"github.com/ramabmtr/inventario/api"
	"github.com/ramabmtr/inventario/config"
	"github.com/ramabmtr/inventario/repository/logger"
)

func init() {
	config.InitEnvVar()
	config.InitLogger()
}

func main() {
	db, err := config.InitDatabaseClient()
	if err != nil {
		logger.WithError(err).Fatal("fail to initialize database connection")
	}
	defer db.Close()

	api.Run()
}

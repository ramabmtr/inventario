// +build main

package main

import (
	"github.com/ramabmtr/inventario/api"
	"github.com/ramabmtr/inventario/config"
)

func main() {
	config.InitEnvVar()
	config.InitLogger()
	db, err := config.InitDatabaseClient()
	if err != nil {
		config.AppLogger.WithError(err).Fatal("fail to initialize database connection")
	}
	defer db.Close()

	api.Run()
}

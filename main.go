// +build main

package main

import (
	"github.com/ramabmtr/inventario/api"
	"github.com/ramabmtr/inventario/config"
)

func main() {
	config.InitEnvVar()
	config.InitLogger()

	api.Run()
}

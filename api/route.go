package api

import (
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/api/handler/healthcheck"
	"github.com/ramabmtr/inventario/api/handler/inventory"
)

func PublicRouter(g *echo.Group) {
	g.GET(
		"/ping",
		healthcheck.Ping,
	)

	g.POST(
		"/inventory",
		inventory.Create,
	)
}

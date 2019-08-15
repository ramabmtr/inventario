package api

import (
	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/api/handler/healthcheck"
)

func PublicRouter(g *echo.Group) {
	g.GET(
		"/ping",
		healthcheck.Ping,
	)
}

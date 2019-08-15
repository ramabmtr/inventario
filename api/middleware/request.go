package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/labstack/echo"
	"github.com/ramabmtr/inventario/config"
)

func RequestLogger() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) (err error) {

			req := c.Request()
			res := c.Response()

			config.AppLogger.Info(fmt.Sprintf("Request started: %s %s", req.Method, req.URL))

			if config.Env.App.Debug {
				var bodyBytes []byte

				if c.Request().Body != nil {
					bodyBytes, _ = ioutil.ReadAll(c.Request().Body)
				}

				c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))

				var body interface{}

				if err = json.Unmarshal(bodyBytes, &body); err != nil {
					body = nil
				}

				bodyBytes, _ = json.Marshal(body)

				config.AppLogger.
					WithField("headers", req.Header).
					WithField("req_body", string(bodyBytes)).
					Debug("Request detail")
			}

			err = next(c)

			config.AppLogger.
				WithField("status", res.Status).
				Info(fmt.Sprintf("Request finished: %s %s", req.Method, req.URL))

			return err
		}
	}
}

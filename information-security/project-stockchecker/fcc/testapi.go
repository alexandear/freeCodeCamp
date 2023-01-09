package fcc

import (
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

func FCC() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			defer func() {
				if c.Request().Method == http.MethodGet && c.Request().URL.Path == "/_api/app-info" {
					var header struct {
						Headers map[string]string `json:"headers"`
					}
					header.Headers = make(map[string]string)
					for name, values := range c.Response().Header() {
						if (len(values)) > 0 {
							header.Headers[strings.ToLower(name)] = values[0]
						}
					}
					c.JSON(http.StatusOK, header)
					return
				}
			}()

			return next(c)
		}
	}
}

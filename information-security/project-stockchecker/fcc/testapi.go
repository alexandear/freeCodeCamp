package fcc

import (
	"net/http"
	"strings"

	"stockchecker/gotest"

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

				if c.Request().Method == http.MethodGet && c.Request().URL.Path == "/_api/get-tests" {
					type result struct {
						Title string `json:"title"`
						State string `json:"state"`
					}
					res, err := gotest.Run(c.Request().Context(), ".", nil, true)
					if err != nil {
						c.Error(err)
						return
					}

					results := make([]result, 0, len(res.TestResults))
					for title, res := range res.TestResults {
						results = append(results, result{
							Title: title,
							State: string(res.Status),
						})
					}
					c.JSON(http.StatusOK, results)
					return
				}
			}()

			return next(c)
		}
	}
}

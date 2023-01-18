package fcc

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func FCC() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		r := c.Request
		if r.Method == http.MethodGet && r.URL.Path == "/_api/app-info" {
			headers := map[string]string{}
			for name, values := range c.Writer.Header() {
				if (len(values)) > 0 {
					headers[strings.ToLower(name)] = values[0]
				}
			}

			var header struct {
				Headers map[string]string `json:"headers"`
			}
			header.Headers = headers

			c.JSON(http.StatusOK, header)
		}
	}
}

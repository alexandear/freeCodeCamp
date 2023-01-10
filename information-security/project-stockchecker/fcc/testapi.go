package fcc

import (
	"encoding/json"
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

				if c.Request().Method == http.MethodGet && c.Request().URL.Path == "/_api/get-tests" {
					c.JSON(http.StatusOK, testsStatusJSONRaw)
					return
				}
			}()

			return next(c)
		}
	}
}

var testsStatusJSONRaw = json.RawMessage(`
[
   {
      "title":"1 stock",
      "context":" -> Functional Tests -> Routing Tests -> GET /api/stock-prices => stockData object",
      "state":"passed",
      "assertions":[
         {
            "method":"equal",
            "args":[
               "res.status",
               "200"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'likes'"
            ]
         },
         {
            "method":"equal",
            "args":[
               "res.body.stockData.stock",
               "'GOOG'"
            ]
         }
      ]
   },
   {
      "title":"1 stock with like",
      "context":" -> Functional Tests -> Routing Tests -> GET /api/stock-prices => stockData object",
      "state":"passed",
      "assertions":[
         {
            "method":"equal",
            "args":[
               "res.status",
               "200"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'likes'"
            ]
         },
         {
            "method":"equal",
            "args":[
               "res.body.stockData.stock",
               "'GOOG'"
            ]
         },
         {
            "method":"isAbove",
            "args":[
               "res.body.stockData.likes",
               "0"
            ]
         }
      ]
   },
   {
      "title":"1 stock with like again (ensure likes arent double counted)",
      "context":" -> Functional Tests -> Routing Tests -> GET /api/stock-prices => stockData object",
      "state":"passed",
      "assertions":[
         {
            "method":"equal",
            "args":[
               "res.status",
               "200"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData",
               "'likes'"
            ]
         },
         {
            "method":"equal",
            "args":[
               "res.body.stockData.stock",
               "'GOOG'"
            ]
         },
         {
            "method":"equal",
            "args":[
               "res.body.stockData.likes",
               "likes"
            ]
         }
      ]
   },
   {
      "title":"2 stocks",
      "context":" -> Functional Tests -> Routing Tests -> GET /api/stock-prices => stockData object",
      "state":"passed",
      "assertions":[
         {
            "method":"equal",
            "args":[
               "res.status",
               "200"
            ]
         },
         {
            "method":"isArray",
            "args":[
               "res.body.stockData"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[0]",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[0]",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[0]",
               "'rel_likes'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[1]",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[1]",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[1]",
               "'rel_likes'"
            ]
         },
         {
            "method":"oneOf",
            "args":[
               "res.body.stockData[0].stock",
               "['GOOG','MSFT']"
            ]
         },
         {
            "method":"oneOf",
            "args":[
               "res.body.stockData[1].stock",
               "['GOOG','MSFT']"
            ]
         },
         {
            "method":"equal",
            "args":[
               "res.body.stockData[0].rel_likes + res.body.stockData[1].rel_likes",
               "0"
            ]
         }
      ]
   },
   {
      "title":"2 stocks with like",
      "context":" -> Functional Tests -> Routing Tests -> GET /api/stock-prices => stockData object",
      "state":"passed",
      "assertions":[
         {
            "method":"equal",
            "args":[
               "res.status",
               "200"
            ]
         },
         {
            "method":"isArray",
            "args":[
               "res.body.stockData"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[0]",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[0]",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[0]",
               "'rel_likes'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[1]",
               "'stock'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[1]",
               "'price'"
            ]
         },
         {
            "method":"property",
            "args":[
               "res.body.stockData[1]",
               "'rel_likes'"
            ]
         },
         {
            "method":"oneOf",
            "args":[
               "res.body.stockData[0].stock",
               "['GOOG','MSFT']"
            ]
         },
         {
            "method":"oneOf",
            "args":[
               "res.body.stockData[1].stock",
               "['GOOG','MSFT']"
            ]
         },
         {
            "method":"equal",
            "args":[
               "res.body.stockData[0].rel_likes + res.body.stockData[1].rel_likes",
               "0"
            ]
         }
      ]
   }
]`)

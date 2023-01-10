package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type StockPriceResp struct {
	StockData StockDataResp `json:"stockData"`
}

type StockDataResp struct {
	Stock string  `json:"stock"`
	Price float64 `json:"price"`
	Likes int     `json:"likes"`
}

type StockDatasResp struct {
	Stock    string  `json:"stock"`
	Price    float64 `json:"price"`
	RelLikes int     `json:"rel_likes"`
}

type StockPricesResp struct {
	StockData []StockDatasResp `json:"stockData"`
}

type Handler struct {
	e         *echo.Echo
	stockServ *StockService
}

func NewHandler(e *echo.Echo, stockServ *StockService) *Handler {
	h := &Handler{
		e:         e,
		stockServ: stockServ,
	}
	api := e.Group("/api")

	api.GET("/stock-prices*", h.StockPrice)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.e.ServeHTTP(w, req)
}

func (h *Handler) StockPrice(ctx echo.Context) error {
	qparams := ctx.QueryParams()
	stocks, ok := qparams["stock"]
	if !ok || len(stocks) < 1 {
		return fmt.Errorf("stock is required")
	}
	ifLike := ctx.QueryParam("like") == "true"
	remoteAddr := ctx.Request().RemoteAddr

	if len(stocks) == 1 {
		param := StockDataParam{
			Stock:      stocks[0],
			IfLike:     ifLike,
			RemoteAddr: remoteAddr,
		}
		sd, err := h.stockServ.StockDataAndLike(ctx.Request().Context(), param)
		if err != nil {
			return fmt.Errorf("stock data: %w", err)
		}

		return ctx.JSON(http.StatusOK, StockPriceResp{
			StockData: StockDataResp{
				Stock: param.Stock,
				Price: sd.Price,
				Likes: sd.LikesCount,
			},
		})
	}

	param := StockDataParam2{
		Stock1:     stocks[0],
		Stock2:     stocks[1],
		IfLike:     ifLike,
		RemoteAddr: remoteAddr,
	}
	sds, err := h.stockServ.StockDataAndLike2(ctx.Request().Context(), param)
	if err != nil {
		return fmt.Errorf("stock data 2: %w", err)
	}

	return ctx.JSON(
		http.StatusOK,
		StockPricesResp{
			StockData: []StockDatasResp{
				{
					Stock:    stocks[0],
					Price:    sds[0].Price,
					RelLikes: sds[0].LikesCount - sds[1].LikesCount,
				},
				{
					Stock:    stocks[1],
					Price:    sds[1].Price,
					RelLikes: sds[1].LikesCount - sds[0].LikesCount,
				},
			},
		},
	)
}

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

	api.GET("/stock-prices", h.StockPrice)

	return h
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.e.ServeHTTP(w, req)
}

func (h *Handler) StockPrice(ctx echo.Context) error {
	stock := ctx.QueryParam("stock")
	sd, err := h.stockServ.StockData(ctx.Request().Context(), stock)
	if err != nil {
		return fmt.Errorf("stock data: %w", err)
	}

	return ctx.JSON(http.StatusOK, StockPriceResp{
		StockData: StockDataResp{
			Stock: stock,
			Price: sd.Price,
			Likes: sd.LikesCount,
		},
	})
}
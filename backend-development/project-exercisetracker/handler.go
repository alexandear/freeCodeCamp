package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type handler struct {
	e        *echo.Echo
	userServ *userService
}

func newHandler(e *echo.Echo, userServ *userService) http.Handler {
	e.Use(middleware.CORS())
	e.File("/", "views/index.html")
	e.File("/style.css", "public/style.css")

	h := &handler{
		e:        e,
		userServ: userServ,
	}

	gapi := e.Group("/api")

	guser := gapi.Group("/users")
	guser.POST("", h.CreateUser)
	guser.GET("", h.Users)

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.e.ServeHTTP(w, req)
}

func (h *handler) CreateUser(ctx echo.Context) error {
	username := ctx.FormValue("username")

	u, err := h.userServ.CreateUser(ctx.Request().Context(), username)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return ctx.JSON(http.StatusOK, u)
}

func (h *handler) Users(ctx echo.Context) error {
	users, err := h.userServ.AllUsers(ctx.Request().Context())
	if err != nil {
		return fmt.Errorf("all users: %w", err)
	}

	return ctx.JSON(http.StatusOK, users)
}

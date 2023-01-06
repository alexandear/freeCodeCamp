package main

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

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

	guser.POST("/:id/exercises", h.CreateExercise)

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

func (h *handler) CreateExercise(ctx echo.Context) error {
	userID := ctx.Param("id")
	description := ctx.FormValue("description")
	duration, err := strconv.Atoi(ctx.FormValue("duration"))
	if err != nil {
		return fmt.Errorf("duration invalid: %w", err)
	}
	dateStr := ctx.FormValue("date")
	var date time.Time
	if dateStr != "" {
		date = time.Now().UTC()
		d, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return fmt.Errorf("date parse: %w", err)
		}
		date = d
	} else {
		date = time.Now().UTC()
	}

	ex, err := h.userServ.CreateExercise(ctx.Request().Context(), userID, description, duration, date)
	if err != nil {
		return fmt.Errorf("create exercise: %w", err)
	}

	return ctx.JSON(http.StatusOK, ex)
}

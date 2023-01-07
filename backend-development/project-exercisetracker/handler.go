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

type handlerUser struct {
	ID       string `json:"_id"`
	Username string `json:"username"`
}

type handlerExercise struct {
	ID          string `json:"_id"`
	Username    string `json:"username"`
	Date        string `json:"date"`
	Duration    int    `json:"duration"`
	Description string `json:"description"`
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

	userID, err := h.userServ.CreateUser(ctx.Request().Context(), username)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return ctx.JSON(http.StatusOK, makeHandlerUser(userID, username))
}

func (h *handler) Users(ctx echo.Context) error {
	users, err := h.userServ.AllUsers(ctx.Request().Context())
	if err != nil {
		return fmt.Errorf("all users: %w", err)
	}

	handlerUsers := make([]handlerUser, 0, len(users))
	for _, u := range users {
		handlerUsers = append(handlerUsers, makeHandlerUser(u.ID, u.Username))
	}

	return ctx.JSON(http.StatusOK, handlerUsers)
}

func (h *handler) CreateExercise(ctx echo.Context) error {
	userID := ctx.Param("id")
	description := ctx.FormValue("description")
	duration := ctx.FormValue("duration")
	date := ctx.FormValue("date")

	u, ex, err := h.userServ.CreateExercise(ctx.Request().Context(), userID, description, duration, date)
	if err != nil {
		return fmt.Errorf("create exercise: %w", err)
	}

	return ctx.JSON(http.StatusOK, makeHandlerExercise(u, ex))
}

func makeHandlerUser(id, username string) handlerUser {
	return handlerUser{
		ID:       id,
		Username: username,
	}
}

func makeHandlerExercise(u user, ex exercise) handlerExercise {
	return handlerExercise{
		ID:          u.ID,
		Username:    u.Username,
		Date:        ex.Date.Format("Mon Jan 01 2006"),
		Duration:    int(ex.Duration.Minutes()),
		Description: ex.Description,
	}
}

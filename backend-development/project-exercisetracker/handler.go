package main

import (
	"fmt"
	"net/http"
	"strconv"

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
	Date        string `json:"date"`
	Duration    int    `json:"duration"`
	Description string `json:"description"`
}

type handlerLog struct {
	handlerUser
	Count int               `json:"count"`
	Log   []handlerExercise `json:"log"`
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

	guser.GET("/:id/logs", h.Logs)

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

	return ctx.JSON(http.StatusOK, struct {
		handlerUser
		handlerExercise
	}{
		handlerUser:     makeHandlerUser(u.ID, u.Username),
		handlerExercise: makeHandlerExercise(ex),
	})
}

func (h *handler) Logs(ctx echo.Context) error {
	userID := ctx.Param("id")
	limit, _ := strconv.Atoi(ctx.QueryParam("limit")) // 0 on error is acceptable
	from := ctx.QueryParam("from")
	to := ctx.QueryParam("to")

	u, exercises, err := h.userServ.Logs(ctx.Request().Context(), userID, limit, from, to)
	if err != nil {
		return fmt.Errorf("logs: %w", err)
	}

	return ctx.JSON(http.StatusOK, makeHandlerLog(u, exercises))
}

func makeHandlerUser(id, username string) handlerUser {
	return handlerUser{
		ID:       id,
		Username: username,
	}
}

func makeHandlerExercise(ex exercise) handlerExercise {
	return handlerExercise{
		Date:        ex.Date.Format("Mon Jan 01 2006"),
		Duration:    int(ex.Duration.Minutes()),
		Description: ex.Description,
	}
}

func makeHandlerLog(u user, exercises []exercise) handlerLog {
	handlerExercises := make([]handlerExercise, 0, len(exercises))
	for _, ex := range exercises {
		handlerExercises = append(handlerExercises, makeHandlerExercise(ex))
	}

	return handlerLog{
		handlerUser: makeHandlerUser(u.ID, u.Username),
		Count:       len(exercises),
		Log:         handlerExercises,
	}
}

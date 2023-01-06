package main

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type handler struct {
	e  *echo.Echo
	mg *mongo.Client

	db *mongo.Database
}

func newHandler(e *echo.Echo, mongoClient *mongo.Client) http.Handler {
	e.Use(middleware.CORS())
	e.File("/", "views/index.html")
	e.File("/style.css", "public/style.css")

	h := &handler{
		e:  e,
		mg: mongoClient,

		db: mongoClient.Database("exercise_tracker"),
	}

	gapi := e.Group("/api")

	guser := gapi.Group("/users")
	guser.POST("", h.CreateUser)

	return h
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	h.e.ServeHTTP(w, req)
}

func (h *handler) CreateUser(ctx echo.Context) error {
	username := ctx.FormValue("username")

	users := h.db.Collection("users")
	res, err := users.InsertOne(ctx.Request().Context(), bson.D{{"username", username}})
	if err != nil {
		return fmt.Errorf("insert one: %w", err)
	}

	return ctx.JSON(http.StatusOK, &user{
		ID:       res.InsertedID.(primitive.ObjectID).Hex(),
		Username: username,
	})
}

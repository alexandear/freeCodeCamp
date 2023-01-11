// Package api provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.12.4 DO NOT EDIT.
package api

import (
	"time"
)

// CreateThreadBody defines model for CreateThreadBody.
type CreateThreadBody struct {
	DeletePassword string `json:"delete_password"`
	Text           string `json:"text"`
}

// Thread defines model for Thread.
type Thread struct {
	Id        string    `json:"_id"`
	BumpedOn  time.Time `json:"bumped_on"`
	CreatedOn time.Time `json:"created_on"`
	Replies   []string  `json:"replies"`
	Reported  bool      `json:"reported"`
	Text      string    `json:"text"`
}

// Board defines model for Board.
type Board = string

// CreateThreadJSONRequestBody defines body for CreateThread for application/json ContentType.
type CreateThreadJSONRequestBody = CreateThreadBody

// CreateThreadFormdataRequestBody defines body for CreateThread for application/x-www-form-urlencoded ContentType.
type CreateThreadFormdataRequestBody = CreateThreadBody

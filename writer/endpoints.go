package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tinrab/meower/util"
)

type CreateMeowRequest struct {
	Body string `json:"body"`
}

func createMeowEndpoint(c *gin.Context) {
	var req CreateMeowRequest
	if err := c.BindJSON(&req); err != nil {
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid request parameters")
		return
	}

	id, err := util.GenerateID()
	if err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to meow")
		return
	}

	if err := dbSession.Query(
		"INSERT INTO meows (id, user_name, body) VALUES (?, ?, ?)",
		id, "tinrab", req.Body,
	).Exec(); err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to meow")
		return
	}

	c.JSON(http.StatusCreated, &Meow{
		ID:       id,
		UserName: "tinrab",
		Body:     req.Body,
	})
}

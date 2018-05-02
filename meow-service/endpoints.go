package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/segmentio/ksuid"
	"github.com/tinrab/meower/meow-service/db"
	"github.com/tinrab/meower/util"
)

type InsertMeowRequest struct {
	Body string `json:"body" form:"body" validate:"required,gte=1,lte=140"`
}

type InsertMeowResponse struct {
	ID string `json:"id"`
}

func insertMeowEndpoint(c *gin.Context) {
	// Read parameters
	var req InsertMeowRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		log.Println(err)
		util.ErrorResponse(c, http.StatusBadRequest, "Invalid parameters")
		return
	}

	// Create meow
	meow := db.Meow{
		ID:   ksuid.New().String(),
		Body: req.Body,
	}
	if err := db.InsertMeow(meow); err != nil {
		util.ErrorResponse(c, http.StatusInternalServerError, "Failed to insert meow")
		return
	}

	// Return new meow
	c.JSON(http.StatusOK, InsertMeowResponse{ID: meow.ID})
}

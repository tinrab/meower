package main

import (
	"net/http"

	"github.com/segmentio/ksuid"
	"github.com/tinrab/meower/mq"
	"github.com/tinrab/meower/schema"
	"github.com/tinrab/meower/util"
)

type CreateMeowResponse struct {
	ID string `json:"id"`
}

func createMeowHandler(w http.ResponseWriter, r *http.Request) {
	// Read parameters
	body := r.FormValue("body")
	if len(body) < 1 || len(body) > 140 {
		util.ResponseError(w, http.StatusBadRequest, "Invalid body")
		return
	}

	// Create meow
	meow := schema.Meow{
		ID:   ksuid.New().String(),
		Body: body,
	}
	if err := mq.WriteMeowCreated(meow.ID, meow.Body); err != nil {
		util.ResponseError(w, http.StatusInternalServerError, "Failed to insert meow")
		return
	}

	// Return new meow
	util.ResponseOk(w, CreateMeowResponse{ID: meow.ID})
}

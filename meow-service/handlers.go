package main

import (
	"log"
	"net/http"

	"github.com/tinrab/meower/mq"

	"github.com/segmentio/ksuid"
	"github.com/tinrab/meower/meow-service/db"
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
	meow := db.Meow{
		ID:   ksuid.New().String(),
		Body: body,
	}
	if err := db.InsertMeow(meow); err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Failed to insert meow")
		return
	}

	mq.WriteMeowCreated(meow.ID, meow.Body)

	// Return new meow
	util.ResponseOk(w, CreateMeowResponse{ID: meow.ID})
}

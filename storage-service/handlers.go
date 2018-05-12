package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/tinrab/meower/mq"
	"github.com/tinrab/meower/schema"
	"github.com/tinrab/meower/storage-service/db"
	"github.com/tinrab/meower/util"
)

type CreateMeowResponse struct {
	ID string `json:"id"`
}

func onMeowCreated(m mq.MeowCreatedMessage) {
	meow := schema.Meow{
		ID:   m.ID,
		Body: m.Body,
	}
	if err := db.InsertMeow(meow); err != nil {
		log.Println(err)
	}
}

func listMeowsHandler(w http.ResponseWriter, r *http.Request) {
	// Read parameters
	skip, err := strconv.ParseUint(r.FormValue("skip"), 10, 64)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
		return
	}
	take, err := strconv.ParseUint(r.FormValue("take"), 10, 64)
	if err != nil {
		util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
		return
	}

	meows, err := db.ListMeows(skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch meows")
		return
	}

	util.ResponseOk(w, meows)
}

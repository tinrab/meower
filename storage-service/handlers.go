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
	var err error
	// Read parameters
	skip := uint64(0)
	skipStr := r.FormValue("skip")
	take := uint64(100)
	takeStr := r.FormValue("take")
	if len(skipStr) != 0 {
		skip, err = strconv.ParseUint(skipStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid skip parameter")
			return
		}
	}
	if len(takeStr) != 0 {
		take, err = strconv.ParseUint(takeStr, 10, 64)
		if err != nil {
			util.ResponseError(w, http.StatusBadRequest, "Invalid take parameter")
			return
		}
	}

	meows, err := db.ListMeows(skip, take)
	if err != nil {
		log.Println(err)
		util.ResponseError(w, http.StatusInternalServerError, "Could not fetch meows")
		return
	}

	util.ResponseOk(w, meows)
}

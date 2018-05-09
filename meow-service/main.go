package main

import (
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tinrab/meower/mq"
)

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", createMeowHandler).
		Methods("POST").
		Queries("body", "{body}")
	return
}

func main() {
	queue := mq.NewKafka("kafka:9092")
	mq.SetMessageQueue(queue)
	defer queue.Close()

	// Run HTTP server
	router := newRouter()
	allowAll := handlers.AllowedOrigins([]string{"*"})
	if err := http.ListenAndServe(":3000", handlers.CORS(allowAll)(router)); err != nil {
		log.Fatal(err)
	}
}

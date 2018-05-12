package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tinrab/retry"

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
	var queue mq.MessageQueue
	err := retry.DoSleep(10, 2*time.Second, func(_ int) error {
		kafka := mq.NewKafka([]string{"kafka:9092"})
		err := kafka.UseProducer()
		queue = kafka
		return err
	})
	if err != nil {
		log.Fatal(err)
	}
	mq.SetMessageQueue(queue)
	defer queue.Close()

	// Run HTTP server
	router := newRouter()
	allowAll := handlers.AllowedOrigins([]string{"*"})
	if err := http.ListenAndServe(":3000", handlers.CORS(allowAll)(router)); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"log"
	"net/http"
	"time"

	"github.com/tinrab/meower/storage-service/db"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/tinrab/meower/mq"
	"github.com/tinrab/retry"
)

func newRouter() (router *mux.Router) {
	router = mux.NewRouter()
	router.HandleFunc("/meows", listMeowsHandler).
		Methods("GET").
		Queries("skip", "{skip}", "take", "{take}")
	return
}

func main() {
	// Connect to PostgreSQL
	err := retry.DoSleep(5, 2*time.Second, func(_ int) error {
		repo, err := db.NewPostgres("postgres://meower:123456@postgres/meower?sslmode=disable")
		if err != nil {
			log.Println(err)
			return err
		}
		db.SetRepository(repo)
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// Subscribe to Kafka
	queue := mq.NewKafka("kafka:9092")
	mq.SetMessageQueue(queue)
	defer queue.Close()
	ch, err := mq.ReadMeow()
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for msg := range ch {
			switch m := msg.(type) {
			case *mq.MeowCreatedMessage:
				onMeowCreated(*m)
			}
		}
	}()

	// Run HTTP server
	router := newRouter()
	allowAll := handlers.AllowedOrigins([]string{"*"})
	if err := http.ListenAndServe(":3000", handlers.CORS(allowAll)(router)); err != nil {
		log.Fatal(err)
	}
}

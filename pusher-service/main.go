package main

import (
	"log"
	"net/http"

	"github.com/tinrab/meower/mq"
)

func main() {
	queue := mq.NewKafka("kafka:9092")
	mq.SetMessageQueue(queue)
	defer queue.Close()

	ch, err := mq.ReadMeow()
	if err != nil {
		log.Fatal(err)
	}

	hub := newHub()

	// Sub to Kafka
	go func() {
		for msg := range ch {
			switch m := msg.(type) {
			case *mq.MeowCreatedMessage:
				log.Printf("meow received: '%v'\n", m)
				hub.broadcast(newMeowCreatedMessage(m.ID, m.Body), nil)
			}
		}
	}()

	// Run WebSocket server
	go hub.run()
	http.HandleFunc("/ws", hub.handleWebSocket)
	err = http.ListenAndServe(":3000", nil)
	if err != nil {
		log.Fatal(err)
	}
}

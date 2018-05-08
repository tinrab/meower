package main

import (
	"log"

	"github.com/tinrab/meower/mq"
)

func main() {
	queue, err := mq.NewRabbitMessageQueue("amqp://guest:guest@rabbitmq:5672")
	if err != nil {
		log.Fatal(err)
	}
	mq.SetMessageQueue(queue)
	defer queue.Close()

	ch, err := mq.Read("meow.#")
	if err != nil {
		log.Fatal(err)
	}

	forever := make(chan struct{})
	go func() {
		for msg := range ch {
			switch m := msg.(type) {
			case *mq.MeowCreatedMessage:
				log.Printf("Meow(%s) created: '%s'\n", m.ID, m.Body)
			}
		}
	}()
	<-forever
}

package main

import (
	"log"

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

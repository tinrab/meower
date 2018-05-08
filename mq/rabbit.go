package mq

import (
	"bytes"
	"encoding/gob"

	"github.com/streadway/amqp"
)

type RabbitMessageQueue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func NewRabbitMessageQueue(address string) (mq *RabbitMessageQueue, err error) {
	mq = &RabbitMessageQueue{}
	mq.conn, err = amqp.Dial(address)
	if err != nil {
		return
	}
	mq.ch, err = mq.conn.Channel()
	if err != nil {
		return
	}
	err = mq.ch.ExchangeDeclare(
		"meower",
		"topic",
		true,
		false,
		false,
		false,
		nil,
	)
	return
}

func (mq *RabbitMessageQueue) Close() {
	mq.conn.Close()
	mq.ch.Close()
}

func (mq *RabbitMessageQueue) WriteMeowCreated(id string, body string) error {
	err := mq.writeMessage(MeowCreatedMessage{
		ID:   id,
		Body: body,
	}, "meow.created")
	return err
}

func (mq *RabbitMessageQueue) Read(routingKey string) (<-chan Message, error) {
	q, err := mq.ch.QueueDeclare(
		"",
		false,
		false,
		true,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	err = mq.ch.QueueBind(
		q.Name,
		routingKey,
		"meower",
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}
	msgs, err := mq.ch.Consume(
		q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, err
	}

	recv := make(chan Message)
	go func() {
		meowCreated := &MeowCreatedMessage{}
		for msg := range msgs {
			switch msg.RoutingKey {
			case "meow.created":
				mq.readMessage(msg.Body, meowCreated)
				recv <- meowCreated
			}
		}
		close(recv)
	}()

	return (<-chan Message)(recv), nil
}

func (mq *RabbitMessageQueue) writeMessage(msg Message, key string) error {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(msg)
	if err != nil {
		return err
	}
	err = mq.ch.Publish(
		"meower",
		key,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        b.Bytes(),
		},
	)
	return err
}

func (mq *RabbitMessageQueue) readMessage(data []byte, m Message) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

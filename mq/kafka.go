package mq

import (
	"bytes"
	"context"
	"encoding/gob"
	"time"

	"github.com/segmentio/kafka-go"
)

type KafkaMessageQueue struct {
	address    string
	meowWriter *kafka.Writer
	meowReader *kafka.Reader
}

func NewKafka(address string) *KafkaMessageQueue {
	return &KafkaMessageQueue{
		address: address,
	}
}

func (mq *KafkaMessageQueue) Close() {
	if mq.meowWriter != nil {
		mq.meowWriter.Close()
	}
	if mq.meowReader != nil {
		mq.meowReader.Close()
	}
}

func (mq *KafkaMessageQueue) WriteMeowCreated(id string, body string) error {
	return mq.writeMeow(&MeowCreatedMessage{
		ID:   id,
		Body: body,
	})
}

func (mq *KafkaMessageQueue) writeMeow(m Message) error {
	if mq.meowWriter == nil {
		mq.meowWriter = kafka.NewWriter(kafka.WriterConfig{
			Brokers:      []string{mq.address},
			Topic:        "meow",
			Balancer:     &kafka.LeastBytes{},
			BatchSize:    0,
			BatchTimeout: 0,
		})
	}

	data, err := mq.writeMessage(m)
	if err != nil {
		return err
	}

	return mq.meowWriter.WriteMessages(context.Background(), kafka.Message{
		Key:   []byte(m.Key()),
		Value: data,
	})
}

func (mq *KafkaMessageQueue) ReadMeow() (<-chan Message, error) {
	if mq.meowReader == nil {
		mq.meowReader = kafka.NewReader(kafka.ReaderConfig{
			Brokers:        []string{mq.address},
			Topic:          "meow",
			MinBytes:       64,
			MaxBytes:       1e3,
			CommitInterval: time.Second,
		})
	}

	recv := make(chan Message)
	go func() {
		meowCreated := &MeowCreatedMessage{}
		ctx := context.Background()
		for {
			km, err := mq.meowReader.ReadMessage(ctx)
			if err != nil {
				break
			}
			switch string(km.Key) {
			case meowCreated.Key():
				err = mq.readMessage(km.Value, meowCreated)
				if err == nil {
					recv <- meowCreated
				}
			}
		}
		close(recv)
	}()
	return (<-chan Message)(recv), nil
}

func (mq *KafkaMessageQueue) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (mq *KafkaMessageQueue) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

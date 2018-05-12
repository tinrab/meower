package mq

import (
	"bytes"
	"encoding/gob"
	"errors"
	"log"
	"os"
	"os/signal"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type KafkaMessageQueue struct {
	brokers  []string
	producer sarama.SyncProducer
	consumer *cluster.Consumer
}

func NewKafka(brokers []string) *KafkaMessageQueue {
	return &KafkaMessageQueue{
		brokers: brokers,
	}
}

func (mq *KafkaMessageQueue) UseProducer() error {
	cfg := sarama.NewConfig()
	// Wait for all in-sync replicas to ack the message
	cfg.Producer.RequiredAcks = sarama.WaitForAll
	// Retry up to 10 times to produce the message
	cfg.Producer.Retry.Max = 10
	cfg.Producer.Return.Successes = true
	cfg.Consumer.Return.Errors = true

	var err error
	mq.producer, err = sarama.NewSyncProducer(mq.brokers, cfg)

	return err
}

func (mq *KafkaMessageQueue) UseConsumer(groupID string) error {
	cfg := cluster.NewConfig()
	cfg.Consumer.Return.Errors = true

	var err error
	mq.consumer, err = cluster.NewConsumer(mq.brokers, groupID, []string{"meow"}, cfg)

	return err
}

func (mq *KafkaMessageQueue) Close() {
	if mq.producer != nil {
		mq.producer.Close()
	}
	if mq.consumer != nil {
		mq.consumer.Close()
	}
}

func (mq *KafkaMessageQueue) WriteMeowCreated(id string, body string) error {
	if mq.producer == nil {
		return errors.New("Producer not connected")
	}
	m := &MeowCreatedMessage{
		ID:   id,
		Body: body,
	}
	data, err := mq.writeMessage(m)
	if err != nil {
		return err
	}
	_, _, err = mq.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "meow",
		// Set key to ensure correct delivery order
		Key:   sarama.StringEncoder(m.Key()),
		Value: sarama.ByteEncoder(data),
	})
	return err
}

func (mq *KafkaMessageQueue) ReadMeow() (<-chan Message, error) {
	if mq.consumer == nil {
		return nil, errors.New("Consumer not connected")
	}

	go func() {
		for err := range mq.consumer.Errors() {
			log.Println(err)
		}
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)
	recv := make(chan Message)
	go func() {
		meowCreated := &MeowCreatedMessage{}
		for {
			select {
			case msg, ok := <-mq.consumer.Messages():
				log.Println(msg)

				if ok {
					switch string(msg.Key) {
					case meowCreated.Key():
						// Decode message
						mq.readMessage(msg.Value, meowCreated)
						recv <- meowCreated
						// Mark as processed
						mq.consumer.MarkOffset(msg, "")
					}
				}
			case <-signals:
				return
			}
		}
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

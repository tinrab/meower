package event

import (
	"bytes"
	"encoding/gob"
	"log"

	"github.com/nats-io/nats.go"

	"github.com/tinrab/meower/schema"
)

type NatsEventStore struct {
	nc                      *nats.Conn
	meowCreatedSubscription *nats.Subscription
	meowCreatedChan         chan MeowCreatedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	nc, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}
	return &NatsEventStore{nc: nc}, nil
}

func (es *NatsEventStore) SubscribeMeowCreated() (<-chan MeowCreatedMessage, error) {
	m := MeowCreatedMessage{}
	es.meowCreatedChan = make(chan MeowCreatedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	es.meowCreatedSubscription, err = es.nc.ChanSubscribe(m.Key(), ch)
	if err != nil {
		return nil, err
	}
	// Decode message
	go func() {
		for {
			select {
			case msg := <-ch:
				if err := es.readMessage(msg.Data, &m); err != nil {
					log.Fatal(err)
				}
				es.meowCreatedChan <- m
			}
		}
	}()
	return (<-chan MeowCreatedMessage)(es.meowCreatedChan), nil
}

func (es *NatsEventStore) OnMeowCreated(f func(MeowCreatedMessage)) (err error) {
	m := MeowCreatedMessage{}
	es.meowCreatedSubscription, err = es.nc.Subscribe(m.Key(), func(msg *nats.Msg) {
		if err := es.readMessage(msg.Data, &m); err != nil {
			log.Fatal(err)
		}
		f(m)
	})
	return
}

func (es *NatsEventStore) Close() {
	if es.nc != nil {
		es.nc.Close()
	}
	if es.meowCreatedSubscription != nil {
		if err := es.meowCreatedSubscription.Unsubscribe(); err != nil {
			log.Fatal(err)
		}
	}
	close(es.meowCreatedChan)
}

func (es *NatsEventStore) PublishMeowCreated(meow schema.Meow) error {
	m := MeowCreatedMessage{meow.ID, meow.Body, meow.CreatedAt}
	data, err := es.writeMessage(&m)
	if err != nil {
		return err
	}
	return es.nc.Publish(m.Key(), data)
}

func (es *NatsEventStore) writeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}

func (es *NatsEventStore) readMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)
	return gob.NewDecoder(&b).Decode(m)
}

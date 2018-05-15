package event

import (
	"time"
)

type Message interface {
	Key() string
}

type MeowCreatedMessage struct {
	ID        string
	Body      string
	CreatedAt time.Time
}

func (m *MeowCreatedMessage) Key() string {
	return "meow.created"
}

package event

type Message interface {
	Key() string
}

type MeowCreatedMessage struct {
	ID   string
	Body string
}

func (m *MeowCreatedMessage) Key() string {
	return "meow.created"
}

package mq

type MessageQueue interface {
	Close()
	WriteMeowCreated(id string, body string) error
	Read(routingKey string) (<-chan Message, error)
}

var impl MessageQueue

func SetMessageQueue(mq MessageQueue) {
	impl = mq
}

func WriteMeowCreated(id string, body string) error {
	return impl.WriteMeowCreated(id, body)
}

func Read(routingKey string) (<-chan Message, error) {
	return impl.Read(routingKey)
}
